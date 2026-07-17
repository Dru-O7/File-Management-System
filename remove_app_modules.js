const fs = require('fs');

let app = fs.readFileSync('frontend/src/app/app.component.html', 'utf8');

// desktop: remove everything from <!-- Section: Modules / Document Folders --> until <!-- Section: User Workspace -->
app = app.replace(/<!-- Section: Modules \/ Document Folders -->[\s\S]*?(?=<!-- Section: User Workspace -->)/, '');

// mobile: remove everything from the mobile documentTypes div until the next div that has "Submit Document"
app = app.replace(/<div\s*class="pt-4 border-t border-\[var\(--border\)\] mt-4"\s*\*ngIf="documentTypes && documentTypes\.length > 0"[\s\S]*?<\/div>\s*(?=<div class="pt-4 border-t border-\[var\(--border\)\] mt-4">\s*<button\s*class="w-full)/, '');

fs.writeFileSync('frontend/src/app/app.component.html', app);
console.log("App component fixed");
