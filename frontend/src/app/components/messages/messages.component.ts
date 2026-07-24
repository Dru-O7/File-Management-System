import { Component, OnInit, ElementRef, ViewChild, AfterViewChecked } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';

export interface ChatThread {
  contact: {
    id: string;
    name: string;
    email: string;
    role: string;
  };
  lastMessage: any;
  unreadCount: number;
  messages: any[];
}

@Component({
  selector: 'app-messages',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './messages.component.html',
  styleUrls: ['./messages.component.css']
})
export class MessagesComponent implements OnInit, AfterViewChecked {
  @ViewChild('chatScrollContainer') private chatScrollContainer!: ElementRef;

  isLoading: boolean = false;
  allMessages: any[] = [];
  chatThreads: ChatThread[] = [];
  selectedContact: { id: string; name: string; email: string; role: string } | null = null;
  activeConversationMessages: any[] = [];

  // Chat Input Dock State
  chatInput: string = '';
  chatSubject: string = '';
  showSubjectInput: boolean = false;
  isSending: boolean = false;

  // Search & New Chat Modal State
  searchQuery: string = '';
  searchResults: any[] = [];
  isSearching: boolean = false;
  searchError: string | null = null;
  activeTab: 'chats' | 'search' = 'chats';

  // Toast feedback
  toastMessage: string | null = null;
  toastType: 'success' | 'error' = 'success';

  private shouldScrollToBottom: boolean = false;

  constructor(
    private api: ApiService,
    public auth: AuthService,
    private route: ActivatedRoute
  ) {}

  ngOnInit() {
    this.loadData();
    this.checkQueryParams();
  }

  checkQueryParams() {
    this.route.queryParams.subscribe(params => {
      if (params['recipientId'] || params['userId']) {
        const id = params['recipientId'] || params['userId'];
        this.openChatByUserId(id);
      } else if (params['email']) {
        this.openChatByEmail(params['email']);
      }
    });
  }

  openChatByUserId(userId: string) {
    const existingThread = this.chatThreads.find(t => t.contact.id === userId);
    if (existingThread) {
      this.selectContact(existingThread.contact);
      return;
    }
    this.api.searchUsers(userId).subscribe({
      next: (results) => {
        const found = (results || []).find(u => (u.id || u.ID) === userId);
        if (found) {
          this.startChatWithUser(found);
        }
      }
    });
  }

  openChatByEmail(email: string) {
    this.api.getUserByEmail(email).subscribe({
      next: (user) => {
        if (user) {
          this.startChatWithUser(user);
        }
      }
    });
  }

  ngAfterViewChecked() {
    if (this.shouldScrollToBottom) {
      this.scrollToBottom();
      this.shouldScrollToBottom = false;
    }
  }

  get currentUser(): any {
    return this.auth.getCurrentUser();
  }

  loadData() {
    this.isLoading = true;
    this.api.getInboxMessages().subscribe({
      next: (inbox) => {
        this.api.getSentMessages().subscribe({
          next: (sent) => {
            this.combineAndGroupMessages(inbox || [], sent || []);
            this.isLoading = false;
          },
          error: (err) => {
            console.error('Failed to load sent messages:', err);
            this.isLoading = false;
          }
        });
      },
      error: (err) => {
        console.error('Failed to load inbox messages:', err);
        this.isLoading = false;
      }
    });
  }

  normalizeContact(c: any): { id: string; name: string; email: string; role: string } | null {
    if (!c) return null;
    const id = c.id || c.ID || c.recipient_id || c.sender_id;
    if (!id) return null;
    return {
      id: String(id),
      name: c.name || c.Name || c.recipient_name || c.sender_name || 'User',
      email: c.email || c.Email || c.recipient_email || c.sender_email || '',
      role: c.role || c.Role || c.recipient_role || c.sender_role || 'Member'
    };
  }

  combineAndGroupMessages(inbox: any[], sent: any[]) {
    const currentUserId = this.currentUser?.id || this.currentUser?.ID;
    this.allMessages = [...inbox, ...sent];

    const threadMap = new Map<string, ChatThread>();

    for (const msg of this.allMessages) {
      const isIncoming = msg.recipient_id === currentUserId;
      const otherId = isIncoming ? msg.sender_id : msg.recipient_id;
      const otherName = isIncoming ? msg.sender_name : msg.recipient_name;
      const otherEmail = isIncoming ? msg.sender_email : msg.recipient_email;
      const otherRole = isIncoming ? msg.sender_role : msg.recipient_role;

      if (!otherId) continue;

      if (!threadMap.has(otherId)) {
        threadMap.set(otherId, {
          contact: {
            id: String(otherId),
            name: otherName || 'User',
            email: otherEmail || '',
            role: otherRole || 'Member'
          },
          lastMessage: msg,
          unreadCount: 0,
          messages: []
        });
      }

      const thread = threadMap.get(otherId)!;
      thread.messages.push(msg);

      if (isIncoming && !msg.is_read) {
        thread.unreadCount++;
      }
    }

    // Sort messages in each thread chronologically (oldest to newest for chat layout)
    this.chatThreads = Array.from(threadMap.values()).map(t => {
      t.messages.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
      t.lastMessage = t.messages[t.messages.length - 1];
      return t;
    });

    // Sort threads list by newest message time
    this.chatThreads.sort((a, b) => {
      const timeA = a.lastMessage ? new Date(a.lastMessage.created_at).getTime() : 0;
      const timeB = b.lastMessage ? new Date(b.lastMessage.created_at).getTime() : 0;
      return timeB - timeA;
    });

    // Refresh active conversation if one is selected
    if (this.selectedContact) {
      const activeThread = this.chatThreads.find(t => t.contact.id === this.selectedContact?.id);
      if (activeThread) {
        this.activeConversationMessages = activeThread.messages;
      }
    } else if (this.chatThreads.length > 0) {
      this.selectContact(this.chatThreads[0].contact);
    }
  }

  selectContact(contact: any) {
    const norm = this.normalizeContact(contact);
    if (!norm) return;

    this.selectedContact = norm;
    let thread = this.chatThreads.find(t => t.contact.id === norm.id);

    if (thread) {
      this.activeConversationMessages = thread.messages;

      // Mark unread messages as read
      const currentUserId = this.currentUser?.id || this.currentUser?.ID;
      for (const msg of thread.messages) {
        if (!msg.is_read && msg.recipient_id === currentUserId) {
          this.api.getMessageDetail(msg.id).subscribe({
            next: () => {
              msg.is_read = true;
              if (thread && thread.unreadCount > 0) thread.unreadCount--;
            }
          });
        }
      }
    } else {
      this.activeConversationMessages = [];
    }

    this.shouldScrollToBottom = true;
  }

  onSearchInput() {
    this.searchError = null;
    const q = this.searchQuery.trim();
    if (q.length < 1) {
      this.searchResults = [];
      return;
    }

    this.isSearching = true;
    this.api.searchUsers(q).subscribe({
      next: (results) => {
        const currentUserId = this.currentUser?.id || this.currentUser?.ID;
        this.searchResults = (results || []).filter(u => (u.id || u.ID) !== currentUserId);
        this.isSearching = false;
        if (this.searchResults.length === 0 && q.includes('@')) {
          this.searchExactEmail();
        }
      },
      error: (err) => {
        console.error('User search error:', err);
        this.searchResults = [];
        this.isSearching = false;
      }
    });
  }

  searchExactEmail() {
    const email = this.searchQuery.trim();
    if (!email) return;

    this.isSearching = true;
    this.searchError = null;

    this.api.getUserByEmail(email).subscribe({
      next: (user) => {
        this.isSearching = false;
        const currentId = this.currentUser?.id || this.currentUser?.ID;
        const userId = user?.id || user?.ID;
        if (user && userId === currentId) {
          this.searchError = 'You cannot chat with yourself.';
          return;
        }
        if (user) {
          this.startChatWithUser(user);
        }
      },
      error: (err) => {
        this.isSearching = false;
        this.searchError = err.error?.error || 'No registered user found with this exact email.';
      }
    });
  }

  startChatWithUser(user: any) {
    const norm = this.normalizeContact(user);
    if (!norm) {
      this.showToast('Unable to start chat with invalid contact.', 'error');
      return;
    }

    const currentUserId = this.currentUser?.id || this.currentUser?.ID;
    if (norm.id === currentUserId) {
      this.searchError = 'You cannot chat with yourself.';
      return;
    }

    let existingThread = this.chatThreads.find(t => t.contact.id === norm.id);
    if (!existingThread) {
      existingThread = {
        contact: norm,
        lastMessage: null,
        unreadCount: 0,
        messages: []
      };
      this.chatThreads.unshift(existingThread);
    }

    this.searchQuery = '';
    this.searchResults = [];
    this.searchError = null;
    this.activeTab = 'chats';
    this.selectContact(norm);
  }

  sendMessage() {
    if (!this.selectedContact) {
      this.showToast('Please select a contact to message.', 'error');
      return;
    }

    const recipientId = this.selectedContact.id;
    if (!recipientId) {
      this.showToast('Invalid recipient selected.', 'error');
      return;
    }

    if (!this.chatInput.trim()) return;

    const body = this.chatInput.trim();
    const subject = this.chatSubject.trim() || 'Chat Message';

    this.isSending = true;
    this.api.sendMessage(recipientId, subject, body).subscribe({
      next: (res) => {
        this.isSending = false;
        this.chatInput = '';
        this.chatSubject = '';
        this.showSubjectInput = false;

        // Find existing thread or create one
        let thread = this.chatThreads.find(t => t.contact.id === recipientId);
        if (!thread) {
          thread = {
            contact: this.selectedContact!,
            lastMessage: res,
            unreadCount: 0,
            messages: []
          };
          this.chatThreads.unshift(thread);
        }

        // Push message to thread messages and update last message
        thread.messages.push(res);
        thread.lastMessage = res;

        // Keep activeConversationMessages in sync
        this.activeConversationMessages = thread.messages;

        // Re-sort thread list by latest message time
        this.chatThreads.sort((a, b) => {
          const timeA = a.lastMessage ? new Date(a.lastMessage.created_at).getTime() : 0;
          const timeB = b.lastMessage ? new Date(b.lastMessage.created_at).getTime() : 0;
          return timeB - timeA;
        });

        this.shouldScrollToBottom = true;
      },
      error: (err) => {
        this.isSending = false;
        this.showToast(err.error?.error || 'Failed to send message.', 'error');
      }
    });
  }

  toggleSubjectInput() {
    this.showSubjectInput = !this.showSubjectInput;
  }

  get totalUnreadCount(): number {
    return this.chatThreads.reduce((sum, t) => sum + t.unreadCount, 0);
  }

  private scrollToBottom(): void {
    try {
      if (this.chatScrollContainer) {
        this.chatScrollContainer.nativeElement.scrollTop = this.chatScrollContainer.nativeElement.scrollHeight;
      }
    } catch (err) {
      console.error('Scroll to bottom error:', err);
    }
  }

  private showToast(msg: string, type: 'success' | 'error' = 'success') {
    this.toastMessage = msg;
    this.toastType = type;
    setTimeout(() => {
      this.toastMessage = null;
    }, 4000);
  }
}

