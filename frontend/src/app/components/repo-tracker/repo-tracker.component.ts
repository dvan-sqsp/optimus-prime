import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Repo } from '../../models/repo.model';
import { RepoService } from '../../services/repo.service';
import { PRListComponent } from '../pull-request-list/pull-request-list.component'
import { ClickOutsideDirective } from '../../directives/click-outside.directive';

@Component({
  selector: 'app-repo-tracker',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, PRListComponent, ClickOutsideDirective],
  templateUrl: './repo-tracker.component.html',
  styleUrls: ['./repo-tracker.component.scss']
})
export class RepoTrackerComponent implements OnInit {
  repoForm: FormGroup;
  repos: Repo[] = [];
  loading = false;
  error = '';

  // Add this property to track which menu is open
  openMenuId: string | null = null;

  // Add these properties
  selectedRepoOwner: string | null = null;
  selectedRepoName: string | null = null;

  // Store the position for the menu
  menuPosition = { x: 0, y: 0 };

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
    this.openMenuId = null; // Close menu after action

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

  // Toggle menu method with improved positioning
  toggleMenu(id: string, event: MouseEvent, index: number): void {
    event.stopPropagation();

    if (this.openMenuId === id) {
      this.closeMenu();
      return;
    }

    // Calculate position that ensures the menu is fully visible
    const buttonRect = (event.target as Element).getBoundingClientRect();

    // Position the menu to the left of the button to avoid right edge clipping
    const menuWidth = 192; // w-48 = 12rem = 192px

    // Default position (show below and to the left of the dots)
    let xPos = buttonRect.left - menuWidth + 20;
    let yPos = buttonRect.bottom + 5;

    // If menu would go off the left edge, position it differently
    if (xPos < 10) {
      xPos = buttonRect.left;
    }

    // If menu would go off the bottom of the viewport, show it above the button
    if (yPos + 100 > window.innerHeight) {
      yPos = buttonRect.top - 100;
    }

    this.menuPosition = { x: xPos, y: yPos };
    this.openMenuId = id;
  }

  // Close menu method
  closeMenu(): void {
    this.openMenuId = null;
  }
}
