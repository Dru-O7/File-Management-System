import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-central-repository',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './central-repository.component.html',
  styleUrls: ['./central-repository.component.css']
})
export class CentralRepositoryComponent implements OnInit {
  loading = true;
  files: any[] = [];
  filteredFiles: any[] = [];
  pendingRequests: any[] = [];
  currentUser: any = null;
  searchText = '';
  activeView: 'repository' | 'requests' = 'repository';

  // Request Access Modal state
  showRequestModal = false;
  selectedFile: any = null;
  requestRemarks = '';
  requestError = '';
  requestSuccess = '';

  // Approve Request Modal state
  showApproveModal = false;
  selectedRequest: any = null;
  approveDuration = 2; // Default 2 hours
  approveError = '';

  // Stat Counters
  totalClosed = 0;
  totalArchived = 0;
  totalPending = 0;

  constructor(
    private api: ApiService,
    private auth: AuthService,
    public router: Router
  ) {}

  ngOnInit() {
    this.currentUser = this.auth.getCurrentUser();
    if (!this.currentUser) {
      this.router.navigate(['/login']);
      return;
    }
    this.loadData();
  }

  loadData() {
    this.loading = true;
    this.api.getClosedOrArchivedFiles(this.searchText).subscribe({
      next: (res) => {
        this.files = res || [];
        this.filteredFiles = this.files;
        this.updateCounters();
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load repository files:', err);
        this.loading = false;
      }
    });

    if (this.isAuthority()) {
      this.loadPendingRequests();
    }
  }

  loadPendingRequests() {
    this.api.getPendingAccessRequests().subscribe({
      next: (res) => {
        this.pendingRequests = res || [];
        this.totalPending = this.pendingRequests.length;
      },
      error: (err) => {
        console.error('Failed to load pending access requests:', err);
      }
    });
  }

  updateCounters() {
    this.totalClosed = this.files.filter(f => f.Status === 'Closed').length;
    this.totalArchived = this.files.filter(f => f.Status === 'Archived').length;
  }

  isAuthority(): boolean {
    if (!this.currentUser) return false;
    const role = this.currentUser.Role || this.currentUser.role;
    return role === 'SuperAdmin' || role === 'Admin' || role === 'School Admin' || role === 'Principal' || role === 'DHE' || (role && role.startsWith('Admin '));
  }

  onSearch() {
    this.loadData();
  }

  openRequestModal(file: any) {
    this.selectedFile = file;
    this.requestRemarks = '';
    this.requestError = '';
    this.requestSuccess = '';
    this.showRequestModal = true;
  }

  submitAccessRequest() {
    if (!this.requestRemarks.trim()) {
      this.requestError = 'Please provide a reason for requesting access.';
      return;
    }

    this.api.requestFileAccess(this.selectedFile.ID, this.requestRemarks).subscribe({
      next: () => {
        this.requestSuccess = 'Access request submitted successfully to authorities.';
        this.requestError = '';
        setTimeout(() => {
          this.showRequestModal = false;
          this.loadData();
        }, 1500);
      },
      error: (err) => {
        this.requestError = err.error?.error || 'Failed to submit request. Please try again.';
      }
    });
  }

  openApproveModal(req: any) {
    this.selectedRequest = req;
    this.approveDuration = 2; // Reset default to 2 hours
    this.approveError = '';
    this.showApproveModal = true;
  }

  resolveRequest(status: 'approved' | 'rejected') {
    const hours = status === 'approved' ? Number(this.approveDuration) : 0;
    this.api.resolveAccessRequest(this.selectedRequest.ID, status, hours).subscribe({
      next: () => {
        this.showApproveModal = false;
        this.loadPendingRequests();
        this.loadData();
      },
      error: (err) => {
        this.approveError = err.error?.error || 'Failed to resolve request.';
      }
    });
  }

  getExpiryLabel(file: any): string {
    // If not a standard user access model or doesn't have an expiry
    if (file.HasAccess) return 'Available';
    return 'No Access';
  }

  goToDetails(file: any) {
    if (file.HasAccess) {
      this.router.navigate(['/details', file.ID], { queryParams: { type: 'file', source: 'repo' } });
    }
  }
}
