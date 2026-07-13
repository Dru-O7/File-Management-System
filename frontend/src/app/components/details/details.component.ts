import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-details',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.css']
})
export class DetailsComponent implements OnInit {
  document: any = null;
  history: any[] = [];
  currentUser: any = null;
  
  actionRemarks: string = '';
  selectedUser: string = '';
  users: any[] = [];
  documentTypes: any[] = [];
  
  selectedFile: File | null = null;
  replaceError: string = '';
  replaceRemarks: string = '';

  pdfCacheBuster: number = Date.now();
  safePdfUrl: SafeResourceUrl | null = null;
  showForwardSelect: boolean = false;
  loading: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private api: ApiService,
    private auth: AuthService,
    public router: Router,
    private sanitizer: DomSanitizer
  ) {}

  toggleForwardSelect() {
    this.showForwardSelect = !this.showForwardSelect;
  }

  ngOnInit() {
    this.currentUser = this.auth.getCurrentUser();
    if (!this.currentUser) {
      this.router.navigate(['/login']);
      return;
    }

    this.api.getUsers().subscribe({
      next: (res) => {
        const currentId = this.currentUser?.ID || this.currentUser?.id;
        this.users = res.filter(u => (u.id || u.ID) !== currentId);
        if (this.users.length > 0) {
          this.selectedUser = this.users[0].id || this.users[0].ID;
        }
      },
      error: (err) => console.error('Failed to load users:', err)
    });

    this.api.getDocumentTypes().subscribe({
      next: (types) => {
        this.documentTypes = types || [];
      },
      error: (err) => console.error('Failed to load document types:', err)
    });

    this.route.paramMap.subscribe(params => {
      const id = params.get('id');
      if (id) {
        this.loadDetails(id);
      }
    });
  }

  loadDetails(id: string) {
    this.loading = true;
    this.api.getDocumentDetails(id).subscribe({
      next: (res) => {
        this.document = res.document;
        this.history = res.history;
        this.pdfCacheBuster = Date.now();
        
        const token = this.auth.getToken();
        const url = `http://localhost:8080/api/documents/${this.document.ID}/download?token=${token}&cb=${this.pdfCacheBuster}`;
        this.safePdfUrl = this.sanitizer.bypassSecurityTrustResourceUrl(url);

        this.draftContent = res.document.DraftSpace || '';
        
        if (this.isDocx(this.document.Filename)) {
          setTimeout(() => {
            this.renderDocxPreview();
          }, 100);
        }
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load document details:', err);
        this.loading = false;
      }
    });
  }

  download() {
    const token = this.auth.getToken();
    window.open(`http://localhost:8080/api/documents/${this.document.ID}/download?token=${token}`, '_blank');
  }

  submitAction(action: string) {
    this.executeSubmitAction(action, '');
  }

  executeSubmitAction(action: string, signature: string) {
    if ((action === 'Sent Back' || action === 'Rejected') && !this.actionRemarks.trim()) {
      alert(`Please enter your Remarks / Noting Sheet comments for this ${action.toLowerCase()} action.`);
      return;
    }

    let target = null;
    if (action === 'Sent Back' || action === 'Rejected') {
      target = this.document.UploaderID;
    } else if (action === 'Approved') {
      target = this.currentUser.ID; // or specific user
    } else if (action === 'Forwarded') {
      if (!this.selectedUser) {
        alert('Please select a user to forward this document to.');
        return;
      }
      target = this.selectedUser;
    }

    this.api.submitAction(this.document.ID, {
      actor_id: this.currentUser.ID,
      target_id: target,
      action: action,
      remarks: this.actionRemarks,
      signature: signature
    }).subscribe({
      next: () => {
        this.loadDetails(this.document.ID);
        this.actionRemarks = '';
        this.showForwardSelect = false;
      },
      error: (err) => {
        console.error('Failed to submit action:', err);
        alert(err.error?.message || 'Failed to submit action. Please make sure all required fields are filled.');
      }
    });
  }

  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
  }

  isPdf(filename: string): boolean {
    return filename ? filename.toLowerCase().endsWith('.pdf') : false;
  }

  isDocx(filename: string): boolean {
    return filename ? filename.toLowerCase().endsWith('.docx') : false;
  }

  renderDocxPreview() {
    if (!this.document) return;
    const token = this.auth.getToken();
    const url = `http://localhost:8080/api/documents/${this.document.ID}/download?token=${token}&cb=${this.pdfCacheBuster}`;
    
    fetch(url)
      .then(response => response.blob())
      .then(blob => {
        const container = document.getElementById('docx-container');
        if (container) {
          container.innerHTML = '';
          import('docx-preview').then(docx => {
            docx.renderAsync(blob, container, undefined, {
              className: 'docx-rendered',
              inWrapper: true,
              ignoreWidth: true,
              ignoreHeight: true,
              ignoreFonts: false,
              breakPages: false,
              debug: false,
              trimXmlDeclaration: true,
              useBase64URL: true,
              renderHeaders: false,
              renderFooters: false,
              renderFootnotes: false,
              renderEndnotes: false,
              experimental: false
            }).catch(err => {
              console.error('Docx render error:', err);
              container.innerHTML = `<div class="flex items-center justify-center h-full text-rose-500 font-semibold p-6 text-center border-2 border-dashed border-rose-200 rounded-xl bg-rose-50/50">
                <p>Failed to render preview. The document might be too large or complex for the browser previewer. Please use the download button below to view it natively.</p>
              </div>`;
            });
          });
        }
      })
      .catch(err => {
        console.error('Error fetching docx:', err);
      });
  }

  getPdfUrl(): SafeResourceUrl {
    return this.safePdfUrl || '';
  }

  getSafeSignature(signature: string): any {
    if (!signature) return '';
    return this.sanitizer.bypassSecurityTrustUrl(signature);
  }

  replaceFile() {
    const formData = new FormData();
    if (this.selectedFile) {
      formData.append('file', this.selectedFile);
    }
    formData.append('uploader_id', this.currentUser.ID);
    formData.append('target_owner_id', this.selectedUser);
    formData.append('remarks', this.replaceRemarks);
    formData.append('title', this.document.Title);
    formData.append('description', this.document.Description);
    formData.append('category', this.document.Category);
    formData.append('tags', this.document.Tags);
    formData.append('priority', this.document.Priority);
    formData.append('direction', this.document.Direction);

    this.api.replaceDocument(this.document.ID, formData).subscribe({
      next: () => {
        this.loadDetails(this.document.ID);
        this.selectedFile = null;
        this.replaceRemarks = '';
        this.replaceError = '';
        const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement;
        if (fileInput) fileInput.value = '';
      },
      error: () => {
        this.replaceError = 'Failed to resubmit document.';
      }
    });
  }

  recallDocument() {
    if (confirm('Are you sure you want to recall this document back to your queue?')) {
      this.api.recallDocument(this.document.ID).subscribe({
        next: () => {
          this.loadDetails(this.document.ID);
        },
        error: (err) => {
          alert('Failed to recall document. It may have already been acted on.');
        }
      });
    }
  }
}
