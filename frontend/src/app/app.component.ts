import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {RepoTrackerComponent} from './components/repo-tracker/repo-tracker.component';
import {PRListComponent} from './components/pull-request-list/pull-request-list.component';
import {CommonModule} from "@angular/common";

@Component({
  selector: 'app-root',
  imports: [CommonModule, RouterOutlet, RepoTrackerComponent, PRListComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'frontend';
}
