const fs = require('fs');

// 1. admin.component.ts
let adminTs = fs.readFileSync('frontend/src/app/components/admin/admin.component.ts', 'utf8');
adminTs = adminTs.replace(/roles = \['Student', 'Teacher', 'Principal', 'Admin', 'Parent'\];/, "roles = ['DHE', 'School Admin', 'Teaching staff', 'non-teaching', 'vocational'];");
adminTs = adminTs.replace(/role: 'Student'/g, "role: 'vocational'");
adminTs = adminTs.replace(/NeedsParentCosign/g, "needs_parent_cosign"); // Just in case, though we are removing it
adminTs = adminTs.replace(/needs_parent_cosign: (true|false|dt\.NeedsParentCosign),/g, "");
fs.writeFileSync('frontend/src/app/components/admin/admin.component.ts', adminTs);

// 2. admin.component.html
let adminHtml = fs.readFileSync('frontend/src/app/components/admin/admin.component.html', 'utf8');
adminHtml = adminHtml.replace(/<th>Parent Cosign<\/th>/g, "");
adminHtml = adminHtml.replace(/<td[^>]*>\s*<span[^>]*>\s*{{dt\.NeedsParentCosign \? 'Required' : 'No'}}\s*<\/span>\s*<\/td>/g, "");
adminHtml = adminHtml.replace(/<div class="field-group">\s*<input type="checkbox" id="cosign" \[\(ngModel\)\]="docTypeForm\.needs_parent_cosign"[^>]*>\s*<label for="cosign"[^>]*>Requires Parent Co-sign<\/label>\s*<\/div>/g, "");
fs.writeFileSync('frontend/src/app/components/admin/admin.component.html', adminHtml);

// 3. upload.component.ts
let uploadTs = fs.readFileSync('frontend/src/app/components/upload/upload.component.ts', 'utf8');
uploadTs = uploadTs.replace(/if \(currentUser && \(currentUser\.Role === 'Student' \|\| currentUser\.Role === 'Parent' \|\| currentUser\.role === 'Student' \|\| currentUser\.role === 'Parent'\)\) \{[\s\S]*?\}\s*else \{/g, "if (true) {");
uploadTs = uploadTs.replace(/u\.Role !== 'Student' && u\.role !== 'Student'/g, "true");
fs.writeFileSync('frontend/src/app/components/upload/upload.component.ts', uploadTs);

// 4. details.component.ts
let detailsTs = fs.readFileSync('frontend/src/app/components/details/details.component.ts', 'utf8');
detailsTs = detailsTs.replace(/u\.Role !== 'Student' &&/g, "");
detailsTs = detailsTs.replace(/u\.role !== 'Student',/g, "");
fs.writeFileSync('frontend/src/app/components/details/details.component.ts', detailsTs);

// 5. details.component.html
let detailsHtml = fs.readFileSync('frontend/src/app/components/details/details.component.html', 'utf8');
detailsHtml = detailsHtml.replace(/<!-- Student Assignment Submission Panel -->[\s\S]*?(?=<!-- Action Log\/Audit Trail -->)/, "");
fs.writeFileSync('frontend/src/app/components/details/details.component.html', detailsHtml);

// 6. login.component.ts
let loginTs = fs.readFileSync('frontend/src/app/components/login/login.component.ts', 'utf8');
loginTs = loginTs.replace(/activePortal: string = 'student';/, "activePortal: string = 'employee';");
fs.writeFileSync('frontend/src/app/components/login/login.component.ts', loginTs);

// 7. register.component.ts
let regTs = fs.readFileSync('frontend/src/app/components/register/register.component.ts', 'utf8');
regTs = regTs.replace(/activePortal: string = 'student';/, "activePortal: string = 'employee';");
fs.writeFileSync('frontend/src/app/components/register/register.component.ts', regTs);

console.log("Frontend cleaned.");
