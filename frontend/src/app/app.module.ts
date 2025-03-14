import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AppComponent } from './app.component';
import { RepoTrackerComponent } from './components/repo-tracker/repo-tracker.component';
import { ClickOutsideModule } from 'ng-click-outside';
import { PRListComponent } from './components/pull-request-list/pull-request-list.component';
import { ClickOutsideDirective } from './directives/click-outside.directive';

@NgModule({
  declarations: [
    ClickOutsideDirective
  ],
  imports: [
    BrowserModule,
    ReactiveFormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    AppComponent,
    RepoTrackerComponent,
    PRListComponent,
    ClickOutsideModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
