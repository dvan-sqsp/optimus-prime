import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Repo } from '../../models/repo.model';
import { RepoService } from '../../services/repo.service';
import { PRListComponent } from '../pull-request-list/pull-request-list.component'

@Component({
  selector: 'app-repo-tracker',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, PRListComponent],
  templateUrl: './repo-tracker.component.html',
  styleUrls: ['./repo-tracker.component.scss']
})
export class RepoTrackerComponent implements OnInit {
  repoForm: FormGroup;
  repos: Repo[] = [];
  loading = false;
  error = '';

  // Add these properties
  selectedRepoOwner: string | null = null;
  selectedRepoName: string | null = null;

  constructor(
    private fb: FormBuilder,
    private repoService: RepoService
  ) {
    this.repoForm = this.fb.group({
      owner: ['', [Validators.required]],
      name: ['', [Validators.required]]
    });
  }

  // Add this method
  selectRepo(owner: string, name: string): void {
    this.selectedRepoOwner = owner;
    this.selectedRepoName = name;
  }

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.loading = true;
    this.repoService.list().subscribe({
      next: (res) => {
        this.repos = res.repos;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load repositories';
        this.loading = false;
        console.error(err);
      }
    });
  }

  onSubmit(): void {
    if (this.repoForm.invalid) {
      return;
    }

    const { owner, name } = this.repoForm.value;
    this.loading = true;

    this.repoService.add(owner, name).subscribe({
      next: (repo) => {
        this.repos.push(repo)
        this.repoForm.reset();
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to add repository';
        this.loading = false;
        console.error(err);
      }
    });
  }

  delete(id: string): void {
    this.loading = true;

    this.repoService.delete(id).subscribe({
      next: () => {
        this.repos = this.repos.filter(repo => repo.id !== id);
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to delete repository';
        this.loading = false;
        console.error(err);
      }
    });
  }
}
