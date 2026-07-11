import { Component, OnInit, OnDestroy } from '@angular/core';
import { RouterOutlet, Router } from '@angular/router';
import { AuthService } from './services/auth.service';
import { ApiService } from './services/api.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, OnDestroy {
  notifications: any[] = [];
  showNotificationsDropdown: boolean = false;
  unreadCount: number = 0;
  private intervalId: any;

  constructor(
    public authService: AuthService,
    private api: ApiService,
    private router: Router
  ) {}

  ngOnInit() {
    this.authService.currentUser$.subscribe(user => {
      if (user) {
        this.startNotificationsPolling();
      } else {
        this.stopNotificationsPolling();
      }
    });
  }

  ngOnDestroy() {
    this.stopNotificationsPolling();
  }

  startNotificationsPolling() {
    this.loadNotifications();
    this.intervalId = setInterval(() => {
      this.loadNotifications();
    }, 10000); // Poll every 10 seconds
  }

  stopNotificationsPolling() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
      this.intervalId = null;
    }
    this.notifications = [];
    this.unreadCount = 0;
  }

  loadNotifications() {
    this.api.getNotifications().subscribe({
      next: (notifs) => {
        this.notifications = notifs || [];
        this.unreadCount = this.notifications.filter(n => n.Status === 'pending').length;
      }
    });
  }

  toggleNotifications() {
    this.showNotificationsDropdown = !this.showNotificationsDropdown;
    if (this.showNotificationsDropdown) {
      // Mark all read visually
      this.unreadCount = 0;
    }
  }

  getNotificationText(n: any): string {
    try {
      const payload = JSON.parse(n.Payload);
      if (n.Template === 'action_required') {
        const sender = payload.uploader_name || payload.actor_name || 'Staff';
        return `Action required: "${payload.document_title}" submitted by ${sender}`;
      } else if (n.Template === 'approved') {
        return `Document approved: "${payload.document_title}" approved by ${payload.actor_name}`;
      } else if (n.Template === 'rejected') {
        return `Document rejected: "${payload.document_title}" rejected by ${payload.actor_name}`;
      } else if (n.Template === 'sent_back') {
        return `Document sent back for revision: "${payload.document_title}" by ${payload.actor_name}`;
      } else if (n.Template === 'sla_warning') {
        return payload.message || `SLA Warning: "${payload.document_title}" has breached deadline.`;
      }
      return `Update on Document ID: ${n.DocumentID}`;
    } catch (e) {
      return `New document update event received.`;
    }
  }

  logout() {
    this.stopNotificationsPolling();
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
