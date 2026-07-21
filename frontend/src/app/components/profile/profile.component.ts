import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';
import { HttpClient } from '@angular/common/http';
import { forkJoin, of } from 'rxjs';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  currentUser: any = null;
  private apiUrl = 'http://localhost:8080/api';

  // Avatar
  avatarInitials: string = '';
  avatarColor: string = '#3b82f6';
  avatarPreview: string | null = null;
  avatarFile: File | null = null;

  // Password reset
  currentPassword: string = '';
  newPassword: string = '';
  confirmPassword: string = '';
  showCurrentPassword: boolean = false;
  showNewPassword: boolean = false;
  showConfirmPassword: boolean = false;

  // Theme & Avatar tracker
  currentTheme: string = 'light';
  originalTheme: string = 'light';
  originalAvatar: string | null = null;

  // Status flags
  saving: boolean = false;
  saveSuccess: string = '';
  saveError: string = '';
  profileName: string = '';
  originalName: string = '';

  private avatarColors = [
    '#3b82f6', '#8b5cf6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444'
  ];

  constructor(
    public router: Router,
    private auth: AuthService,
    private http: HttpClient
  ) {}

  ngOnInit(): void {
    this.auth.currentUser$.subscribe(user => {
      if (!user) {
        this.router.navigate(['/login']);
        return;
      }
      this.currentUser = user;
      this.profileName = user.Name || user.name || '';
      this.originalName = this.profileName;
      this.avatarPreview = user.Avatar || user.avatar || null;
      this.originalAvatar = this.avatarPreview;
      const name: string = this.profileName;
      const words = name.trim().split(' ').filter((w: string) => w.length > 0);
      this.avatarInitials = words.length >= 2
        ? (words[0][0] + words[words.length - 1][0]).toUpperCase()
        : name.slice(0, 2).toUpperCase();
      const colorIndex = name.charCodeAt(0) % this.avatarColors.length;
      this.avatarColor = this.avatarColors[colorIndex];
    });

    // Initialize theme preference
    this.currentTheme = localStorage.getItem('theme') || 'light';
    this.originalTheme = this.currentTheme;
  }

  onAvatarFileChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    if (!file.type.startsWith('image/')) {
      return;
    }
    this.avatarFile = file;
    const reader = new FileReader();
    reader.onload = (e) => {
      this.avatarPreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  }

  setTheme(theme: string): void {
    this.currentTheme = theme;
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  saveChanges(): void {
    this.saveError = '';
    this.saveSuccess = '';

    const isNameChanged = this.profileName !== this.originalName;
    const isThemeChanged = this.currentTheme !== this.originalTheme;
    const isAvatarChanged = this.avatarPreview !== this.originalAvatar;
    const isPasswordAttempt = !!(this.currentPassword || this.newPassword || this.confirmPassword);

    if (!isNameChanged && !isThemeChanged && !isAvatarChanged && !isPasswordAttempt) {
      this.saveError = 'No changes detected to save.';
      return;
    }

    if (isPasswordAttempt) {
      if (!this.currentPassword || !this.newPassword || !this.confirmPassword) {
        this.saveError = 'All password fields are required to update password.';
        return;
      }
      if (this.newPassword !== this.confirmPassword) {
        this.saveError = 'New password and confirm password do not match.';
        return;
      }
      if (this.newPassword.length < 8) {
        this.saveError = 'New password must be at least 8 characters.';
        return;
      }
    }

    if (!isNameChanged && !isAvatarChanged && !isPasswordAttempt) {
      // Theme only change
      localStorage.setItem('theme', this.currentTheme);
      this.originalTheme = this.currentTheme;
      this.saveSuccess = 'Theme preference updated successfully.';
      setTimeout(() => this.saveSuccess = '', 4000);
      return;
    }

    this.saving = true;

    const nameReq = isNameChanged
      ? this.http.put(`${this.apiUrl}/profile/name`, { name: this.profileName })
      : of(null);

    const avatarReq = isAvatarChanged
      ? this.http.put(`${this.apiUrl}/profile/avatar`, { avatar: this.avatarPreview })
      : of(null);

    const passwordReq = isPasswordAttempt
      ? this.http.put(`${this.apiUrl}/profile/password`, {
          current_password: this.currentPassword,
          new_password: this.newPassword
        })
      : of(null);

    forkJoin({
      name: nameReq,
      avatar: avatarReq,
      password: passwordReq
    }).subscribe({
      next: (res: any) => {
        this.saving = false;
        
        let msg = 'Profile updated successfully.';
        if (isPasswordAttempt && isAvatarChanged) {
          msg = 'Password and profile photo updated successfully.';
        } else if (isPasswordAttempt) {
          msg = 'Password updated successfully.';
        } else if (isAvatarChanged && isThemeChanged) {
          msg = 'Profile photo and theme preference updated.';
        } else if (isAvatarChanged) {
          msg = 'Profile photo updated successfully.';
        } else if (isThemeChanged) {
          msg = 'Theme preference updated successfully.';
        }
        
        this.saveSuccess = msg;

        // Save Name to currentUser state
        if (isNameChanged) {
          const updatedUser = { 
            ...this.currentUser, 
            Name: this.profileName,
            name: this.profileName
          };
          this.currentUser = updatedUser;
          this.originalName = this.profileName;
          this.auth.setCurrentUser(updatedUser, this.auth.getToken()!);
        }

        // Save Theme to localStorage
        if (isThemeChanged) {
          localStorage.setItem('theme', this.currentTheme);
          this.originalTheme = this.currentTheme;
        }

        // Save Avatar to localStorage & currentUser state
        if (isAvatarChanged) {
          const updatedUser = { 
            ...this.currentUser, 
            Avatar: this.avatarPreview,
            avatar: this.avatarPreview
          };
          this.currentUser = updatedUser;
          this.originalAvatar = this.avatarPreview;
          this.auth.setCurrentUser(updatedUser, this.auth.getToken()!);
        }

        // Reset password fields
        this.currentPassword = '';
        this.newPassword = '';
        this.confirmPassword = '';

        setTimeout(() => this.saveSuccess = '', 4000);
      },
      error: (err: any) => {
        this.saving = false;
        this.saveError = err.error?.error || 'Failed to save changes. Please check current password.';
      }
    });
  }
}
