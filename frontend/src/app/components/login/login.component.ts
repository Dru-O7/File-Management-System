import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  email: string = 'alice@school.edu';
  password: string = 'password';
  error: string = '';
  loading: boolean = false;
  activePortal: string = 'employee';

  constructor(private api: ApiService, private auth: AuthService, private router: Router) {}

  login() {
    this.error = '';
    const emailTrimmed = this.email.trim();
    if (!emailTrimmed || !this.password) {
      this.error = 'All fields are required.';
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(emailTrimmed)) {
      this.error = 'Please enter a valid email address.';
      return;
    }

    this.loading = true;
    this.api.login(emailTrimmed, this.password).subscribe({
      next: (res) => {
        this.loading = false;
        this.auth.setCurrentUser(res.user, res.token);
        const role = res.user.Role || res.user.role;
        if (role === 'Admin' || role === 'SuperAdmin') {
          this.router.navigate(['/admin']);
        } else {
          this.router.navigate(['/dashboard']);
        }
      },
      error: () => {
        this.loading = false;
        this.error = 'Invalid email/password or user not found. (Hint: default password is "password")';
      }
    });
  }
}
