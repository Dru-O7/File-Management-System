const fs = require('fs');

// 1. backend/cmd/api/main.go
let apiGo = fs.readFileSync('backend/cmd/api/main.go', 'utf8');

// Remove ParentChild cleanup loop
apiGo = apiGo.replace(/db\.Session\(&gorm\.Session\{AllowGlobalUpdate: true\}\)\.Delete\(&models\.ParentChild\{\}\)/g, "");
apiGo = apiGo.replace(/NeedsParentCosign:\s*(true|false),/g, "");
fs.writeFileSync('backend/cmd/api/main.go', apiGo);

// 2. backend/cmd/seed/main.go - I'll just rewrite the docTypes array properly
let seedGo = fs.readFileSync('backend/cmd/seed/main.go', 'utf8');
const docTypesRegex = /docTypes := \[\]models\.DocumentType\{[\s\S]*?\}(?=\s*for i := range docTypes)/;
const newDocTypes = `docTypes := []models.DocumentType{
			{
				SchoolID:       school.ID,
				Name:           "Staff Grievance",
				Slug:           "staff-grievance",
				WorkflowStages: \`[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]\`,
				RequiredFields: \`[]\`,
				SlaHours:       72,
			},
			{
				SchoolID:       school.ID,
				Name:           "Infrastructure Issue",
				Slug:           "infrastructure-issue",
				WorkflowStages: \`[{"stage": 1, "role": "School Admin", "label": "School Admin Final approval", "optional": false}]\`,
				RequiredFields: \`["reason", "urgency"]\`,
				SlaHours:       120,
			},
			{
				SchoolID:       school.ID,
				Name:           "Disciplinary Issue",
				Slug:           "disciplinary-issue",
				WorkflowStages: \`[{"stage": 1, "role": "Teaching staff", "label": "Department Head", "optional": false}]\`,
				RequiredFields: \`["event_name", "event_date"]\`,
				SlaHours:       24,
			},
			{
				SchoolID:       school.ID,
				Name:           "Audit Report",
				Slug:           "audit-report",
				WorkflowStages: \`[{"stage": 1, "role": "School Admin", "label": "School Admin Approval", "optional": false}]\`,
				RequiredFields: \`["audit_reason", "percentage"]\`,
				SlaHours:       96,
			},
		}`;
seedGo = seedGo.replace(docTypesRegex, newDocTypes);
fs.writeFileSync('backend/cmd/seed/main.go', seedGo);
console.log("Fixed.");
