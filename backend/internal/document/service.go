package document

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type Service interface {
	Upload(uploaderID, targetOwnerID uuid.UUID, title, description, category, tags string, fileHeader *multipart.FileHeader) (*DocumentResponse, error)
	List(userID uuid.UUID, search string) ([]DocumentResponse, error)
	GetDetails(docID, authenticatedUserID uuid.UUID) (*DocumentDetailsResponse, error)
	GetFilePathForDownload(docID, authenticatedUserID uuid.UUID) (string, error)
	Replace(docID, authenticatedUserID, targetOwnerID uuid.UUID, title, description, category, tags string, fileHeader *multipart.FileHeader, remarks string) (*DocumentResponse, error)
	TakeAction(docID, authenticatedUserID uuid.UUID, req ActionRequest) (*DocumentResponse, error)
}

type service struct {
	repo       Repository
	uploadsDir string
}

func NewService(repo Repository, uploadsDir string) Service {
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}
	return &service{repo: repo, uploadsDir: uploadsDir}
}

func (s *service) Upload(uploaderID, targetOwnerID uuid.UUID, title, description, category, tags string, fileHeader *multipart.FileHeader) (*DocumentResponse, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	uniquePrefix := fmt.Sprintf("%d_", time.Now().Unix())
	safeFilename := uniquePrefix + filepath.Base(fileHeader.Filename)
	destPath := filepath.Join(s.uploadsDir, safeFilename)

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	uniqueNum := fmt.Sprintf("DOC-%d", time.Now().UnixNano()/1e6)

	docID := uuid.New()
	doc := &models.Document{
		ID:             docID,
		Filename:       fileHeader.Filename,
		FilePath:       destPath,
		UploaderID:     uploaderID,
		CurrentOwnerID: targetOwnerID,
		Status:         models.StatusPendingApproval,
		Title:          title,
		Description:    description,
		UniqueNumber:   uniqueNum,
		Tags:           tags,
		Category:       category,
	}

	if err := s.repo.Create(doc); err != nil {
		return nil, err
	}

	history := &models.WorkflowHistory{
		ID:         uuid.New(),
		DocumentID: docID,
		ActorID:    uploaderID,
		TargetID:   &targetOwnerID,
		Action:     models.ActionUploaded,
		Remarks:    "Document submitted for approval",
	}

	if err := s.repo.CreateHistory(history); err != nil {
		log.Printf("Warning: Failed to write upload workflow log: %v", err)
	}

	savedDoc, err := s.repo.GetByID(docID)
	if err != nil {
		return nil, err
	}

	return s.toDocumentResponse(savedDoc), nil
}

func (s *service) List(userID uuid.UUID, search string) ([]DocumentResponse, error) {
	docs, err := s.repo.ListByUser(userID, search)
	if err != nil {
		return nil, err
	}

	responses := make([]DocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = *s.toDocumentResponse(&d)
	}
	return responses, nil
}

func (s *service) GetDetails(docID, authenticatedUserID uuid.UUID) (*DocumentDetailsResponse, error) {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if err := s.authorizeDocAccess(doc, authenticatedUserID); err != nil {
		return nil, err
	}

	histories, err := s.repo.GetHistoryByDocumentID(docID)
	if err != nil {
		return nil, err
	}

	docDto := s.toDocumentResponse(doc)
	historyDtos := make([]HistoryResponse, len(histories))
	for i, h := range histories {
		historyDtos[i] = *s.toHistoryResponse(&h)
	}

	return &DocumentDetailsResponse{
		Document: *docDto,
		History:  historyDtos,
	}, nil
}

func (s *service) GetFilePathForDownload(docID, authenticatedUserID uuid.UUID) (string, error) {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return "", errors.New("document not found")
	}

	if err := s.authorizeDocAccess(doc, authenticatedUserID); err != nil {
		return "", err
	}

	return doc.FilePath, nil
}

func (s *service) Replace(docID, authenticatedUserID, targetOwnerID uuid.UUID, title, description, category, tags string, fileHeader *multipart.FileHeader, remarks string) (*DocumentResponse, error) {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if doc.UploaderID != authenticatedUserID {
		return nil, errors.New("only the original uploader is authorized to replace or resubmit this document")
	}

	if doc.Status != models.StatusSentBack {
		return nil, errors.New("document must be in 'Sent Back' status to be replaced or resubmitted")
	}

	fileReplaced := false
	if fileHeader != nil {
		src, err := fileHeader.Open()
		if err == nil {
			defer src.Close()
			uniquePrefix := fmt.Sprintf("%d_", time.Now().Unix())
			safeFilename := uniquePrefix + filepath.Base(fileHeader.Filename)
			destPath := filepath.Join(s.uploadsDir, safeFilename)

			dst, err := os.Create(destPath)
			if err == nil {
				defer dst.Close()
				if _, err = io.Copy(dst, src); err == nil {
					os.Remove(doc.FilePath)
					doc.Filename = fileHeader.Filename
					doc.FilePath = destPath
					fileReplaced = true
				}
			}
		}
	}

	if title != "" {
		doc.Title = title
	}
	if description != "" {
		doc.Description = description
	}
	if category != "" {
		doc.Category = category
	}
	if tags != "" {
		doc.Tags = tags
	}

	doc.Status = models.StatusPendingApproval
	doc.CurrentOwnerID = targetOwnerID

	if err := s.repo.Save(doc); err != nil {
		return nil, err
	}

	wfAction := models.ActionResubmitted
	if fileReplaced {
		wfAction = models.ActionFileReplaced
	}

	history := &models.WorkflowHistory{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		ActorID:    authenticatedUserID,
		TargetID:   &targetOwnerID,
		Action:     wfAction,
		Remarks:    remarks,
	}
	_ = s.repo.CreateHistory(history)

	updatedDoc, _ := s.repo.GetByID(doc.ID)
	return s.toDocumentResponse(updatedDoc), nil
}

func (s *service) TakeAction(docID, authenticatedUserID uuid.UUID, req ActionRequest) (*DocumentResponse, error) {
	doc, err := s.repo.GetByID(docID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if doc.CurrentOwnerID != authenticatedUserID {
		return nil, errors.New("you are not authorized to act on this document as you are not the current owner")
	}

	var newStatus models.DocumentStatus
	var nextOwnerID uuid.UUID
	wfAction := models.WorkflowAction(req.Action)

	switch wfAction {
	case models.ActionApproved:
		newStatus = models.StatusApproved
		nextOwnerID = authenticatedUserID
	case models.ActionRejected:
		newStatus = models.StatusRejected
		nextOwnerID = doc.UploaderID
	case models.ActionSentBack:
		newStatus = models.StatusSentBack
		nextOwnerID = doc.UploaderID
	case models.ActionForwarded, "Forward":
		newStatus = models.StatusPendingApproval
		if req.TargetID == nil {
			return nil, errors.New("target ID is required to forward this document")
		}
		nextOwnerID = *req.TargetID
		wfAction = models.ActionForwarded
	default:
		return nil, errors.New("invalid action name")
	}

	doc.Status = newStatus
	doc.CurrentOwnerID = nextOwnerID

	if req.Signature != "" {
		existingSigs, _ := s.repo.CountSignatures(doc.ID)
		filePathLower := strings.ToLower(doc.FilePath)
		if strings.HasSuffix(filePathLower, ".pdf") {
			if err := stampSignatureOnPDF(doc.FilePath, req.Signature, existingSigs); err != nil {
				log.Printf("Error overlaying signature on PDF: %v", err)
			}
		} else if strings.HasSuffix(filePathLower, ".docx") {
			if err := stampSignatureOnDocx(doc.FilePath, req.Signature, existingSigs); err != nil {
				log.Printf("Error overlaying signature on DOCX: %v", err)
			}
		}
	}

	if err := s.repo.Save(doc); err != nil {
		return nil, err
	}

	history := &models.WorkflowHistory{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		ActorID:    authenticatedUserID,
		TargetID:   req.TargetID,
		Action:     wfAction,
		Remarks:    req.Remarks,
		Signature:  req.Signature,
	}
	_ = s.repo.CreateHistory(history)

	updatedDoc, _ := s.repo.GetByID(doc.ID)
	return s.toDocumentResponse(updatedDoc), nil
}

func (s *service) authorizeDocAccess(doc *models.Document, userID uuid.UUID) error {
	if doc.UploaderID == userID || doc.CurrentOwnerID == userID {
		return nil
	}

	histories, err := s.repo.GetHistoryByDocumentID(doc.ID)
	if err == nil {
		for _, h := range histories {
			if h.ActorID == userID {
				return nil
			}
		}
	}

	return errors.New("you are not authorized to view or access this document")
}

func (s *service) toDocumentResponse(d *models.Document) *DocumentResponse {
	return &DocumentResponse{
		ID:             d.ID,
		Filename:       d.Filename,
		FilePath:       d.FilePath,
		UploaderID:     d.UploaderID,
		CurrentOwnerID: d.CurrentOwnerID,
		Status:         d.Status,
		Title:          d.Title,
		Description:    d.Description,
		UniqueNumber:   d.UniqueNumber,
		Tags:           d.Tags,
		Category:       d.Category,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
		Uploader:       d.Uploader,
		CurrentOwner:   d.CurrentOwner,
	}
}

func (s *service) toHistoryResponse(h *models.WorkflowHistory) *HistoryResponse {
	return &HistoryResponse{
		ID:         h.ID,
		DocumentID: h.DocumentID,
		ActorID:    h.ActorID,
		TargetID:   h.TargetID,
		Action:     h.Action,
		Remarks:    h.Remarks,
		Signature:  h.Signature,
		CreatedAt:  h.CreatedAt,
		Actor:      h.Actor,
		Target:     h.Target,
	}
}

func stampSignatureOnPDF(pdfPath string, base64Signature string, existingSigCount int) error {
	if base64Signature == "" {
		return nil
	}
	if !strings.HasSuffix(strings.ToLower(pdfPath), ".pdf") {
		return nil
	}
	parts := strings.Split(base64Signature, ",")
	base64Data := parts[len(parts)-1]

	dec, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	tempDir := os.TempDir()
	tempPNG := filepath.Join(tempDir, fmt.Sprintf("sig_temp_%d.png", time.Now().UnixNano()))
	err = os.WriteFile(tempPNG, dec, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempPNG)

	tempOutPDF := pdfPath + ".signed"
	offsetX := -20 - (existingSigCount * 110)
	desc := fmt.Sprintf("scale:0.25, pos:br, off:%d 20", offsetX)
	wm, err := pdfcpu.ParseImageWatermarkDetails(tempPNG, desc, true, types.POINTS)
	if err != nil {
		return err
	}

	err = api.AddWatermarksFile(pdfPath, tempOutPDF, nil, wm, nil)
	if err != nil {
		return err
	}

	err = os.Rename(tempOutPDF, pdfPath)
	if err != nil {
		os.Remove(tempOutPDF)
		return err
	}

	return nil
}

func stampSignatureOnDocx(docxPath string, base64Signature string, existingSigCount int) error {
	if base64Signature == "" {
		return nil
	}

	parts := strings.Split(base64Signature, ",")
	base64Data := parts[len(parts)-1]

	dec, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(docxPath)
	if err != nil {
		return err
	}
	defer r.Close()

	tempOutDocx := docxPath + ".signed"
	out, err := os.Create(tempOutDocx)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	sigID := fmt.Sprintf("rIdSig%d", existingSigCount+1)
	sigFileName := fmt.Sprintf("sig%d.png", existingSigCount+1)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		var content []byte
		if f.Name == "word/document.xml" || f.Name == "word/_rels/document.xml.rels" || f.Name == "[Content_Types].xml" {
			content, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
		}

		if f.Name == "word/document.xml" {
			idx := bytes.LastIndex(content, []byte("</w:body>"))
			if idx == -1 {
				return fmt.Errorf("could not find closing w:body tag in word/document.xml")
			}

			cx := 1828800
			cy := 731520
			xmlInsert := fmt.Sprintf(`
<w:p xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:pPr>
    <w:jc w:val="right"/>
  </w:pPr>
  <w:r>
    <w:rPr>
      <w:sz w:val="20"/>
      <w:b/>
    </w:rPr>
    <w:t>Signed electronically:</w:t>
    <w:br/>
  </w:r>
  <w:r>
    <w:drawing xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
      <wp:inline distT="0" distB="0" distL="0" distR="0">
        <wp:extent cx="%d" cy="%d"/>
        <wp:effectExtent l="0" t="0" r="0" b="0"/>
        <wp:docPr id="1000" name="Signature"/>
        <wp:cNvGraphicFramePr>
          <a:graphicFrameLocks noChangeAspect="1"/>
        </wp:cNvGraphicFramePr>
        <a:graphic>
          <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing">
            <pic:pic>
              <pic:nvPicPr>
                <pic:cNvPr id="1000" name="Signature"/>
                <pic:cNvPicPr/>
              </pic:nvPicPr>
              <pic:blipFill>
                <a:blip r:embed="%s"/>
                <a:stretch>
                  <a:fillRect/>
                </a:stretch>
              </pic:blipFill>
              <pic:spPr>
                <a:xfrm>
                  <a:off x="0" y="0"/>
                  <a:ext cx="%d" cy="%d"/>
                </a:xfrm>
                <a:prstGeom prst="rect">
                  <a:avLst/>
                </a:prstGeom>
              </pic:spPr>
            </pic:pic>
          </a:graphicData>
        </a:graphic>
      </wp:inline>
    </w:drawing>
  </w:r>
</w:p>`, cx, cy, sigID, cx, cy)

			newContent := append(content[:idx], []byte(xmlInsert)...)
			newContent = append(newContent, content[idx:]...)

			fw, err := w.Create(f.Name)
			if err != nil {
				return err
			}
			_, err = fw.Write(newContent)
			if err != nil {
				return err
			}

		} else if f.Name == "word/_rels/document.xml.rels" {
			idx := bytes.LastIndex(content, []byte("</Relationships>"))
			if idx == -1 {
				return fmt.Errorf("could not find closing Relationships tag in word/_rels/document.xml.rels")
			}

			relInsert := fmt.Sprintf(`  <Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/%s"/>
`, sigID, sigFileName)

			newContent := append(content[:idx], []byte(relInsert)...)
			newContent = append(newContent, content[idx:]...)

			fw, err := w.Create(f.Name)
			if err != nil {
				return err
			}
			_, err = fw.Write(newContent)
			if err != nil {
				return err
			}

		} else if f.Name == "[Content_Types].xml" {
			var newContent []byte
			if !bytes.Contains(content, []byte(`Extension="png"`)) {
				idx := bytes.Index(content, []byte("<Types"))
				if idx == -1 {
					return fmt.Errorf("could not find Types tag in [Content_Types].xml")
				}
				closeIdx := bytes.Index(content[idx:], []byte(">"))
				if closeIdx == -1 {
					return fmt.Errorf("could not find end of Types tag in [Content_Types].xml")
				}
				insertPos := idx + closeIdx + 1
				decl := []byte(`
  <Default Extension="png" ContentType="image/png"/>`)
				newContent = append(content[:insertPos], decl...)
				newContent = append(newContent, content[insertPos:]...)
			} else {
				newContent = content
			}

			fw, err := w.Create(f.Name)
			if err != nil {
				return err
			}
			_, err = fw.Write(newContent)
			if err != nil {
				return err
			}

		} else {
			fw, err := w.CreateHeader(&f.FileHeader)
			if err != nil {
				rc.Close()
				return err
			}
			_, err = io.Copy(fw, rc)
			rc.Close()
			if err != nil {
				return err
			}
		}
	}

	mediaPath := fmt.Sprintf("word/media/%s", sigFileName)
	fw, err := w.Create(mediaPath)
	if err != nil {
		return err
	}
	_, err = fw.Write(dec)
	if err != nil {
		return err
	}

	w.Close()
	out.Close()
	r.Close()

	err = os.Rename(tempOutDocx, docxPath)
	if err != nil {
		os.Remove(tempOutDocx)
		return err
	}

	return nil
}
