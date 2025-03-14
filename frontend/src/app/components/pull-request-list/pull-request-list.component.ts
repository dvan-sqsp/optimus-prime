import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PullRequest } from '../../models/pr.model';
import { PRService } from '../../services/pr.service';

@Component({
  selector: 'app-pull-request-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './pull-request-list.component.html',
  styleUrls: ['./pull-request-list.component.scss']
})
export class PRListComponent implements OnChanges {
  @Input() repoOwner: string | null = null;
  @Input() repoName: string | null = null;

  pullRequests: PullRequest[] = [];
  loading = false;
  error = '';

  constructor(private prService: PRService) { }

  ngOnChanges(changes: SimpleChanges): void {
    if ((changes['repoOwner'] || changes['repoName']) && this.repoOwner && this.repoName) {
      this.loadPullRequests();
    }
  }

  loadPullRequests(): void {
    if (!this.repoOwner || !this.repoName) {
      return;
    }

    this.loading = true;
    this.pullRequests = [];
    this.error = '';

    this.prService.list(this.repoOwner, this.repoName).subscribe({
      next: (response) => {
        this.pullRequests = response.pull_requests;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load pull requests';
        this.loading = false;
        console.error(err);
      }
    });
  }

  getStatusClass(status: string): string {
    switch (status.toLowerCase()) {
      case 'open':
        return 'bg-green-100 text-green-800';
      case 'closed':
        return 'bg-red-100 text-red-800';
      case 'merged':
        return 'bg-purple-100 text-purple-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  getBorderClass(createdAtStr: string): string {
    const createdAt = new Date(createdAtStr);
    const now = new Date();
    const diffDays = Math.floor((now.getTime() - createdAt.getTime()) / (1000 * 60 * 60 * 24));

    if (diffDays > 5) {
      return 'border-red-500';
    } else if (diffDays > 2) {
      return 'border-yellow-500';
    } else {
      return 'border-gray-200';
    }
  }
}
