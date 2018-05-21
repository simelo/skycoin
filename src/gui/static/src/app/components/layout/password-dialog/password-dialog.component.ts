import { Component, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef, MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { FormControl, FormGroup } from '@angular/forms';
import { ButtonComponent } from '../button/button.component';
import { parseResponseMessage } from '../../../utils/errors';
import { Subject } from 'rxjs/Subject';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.scss'],
})
export class PasswordDialogComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  passwordSubmit = new Subject<any>();
  disableDismiss = false;

  private subscriptions: ISubscription[] = [];

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    private snackbar: MatSnackBar,
  ) {
    this.data = Object.assign({
      confirm: false,
      description: null,
      title: null,
    }, data || {});
  }

  ngOnInit() {
    this.form = new FormGroup({}, this.validateForm.bind(this));
    this.form.addControl('password', new FormControl(''));
    this.form.addControl('confirm_password', new FormControl(''));

    ['password', 'confirm_password'].forEach(control => {
      this.subscriptions.push(this.form.get(control).valueChanges.subscribe(() => {
        if (this.button.state === 2) {
          this.button.resetState();
        }
      }));
    });

    if (this.data.confirm) {
      this.form.get('confirm_password').enable();
    } else {
      this.form.get('confirm_password').disable();
    }

    if (this.data.description) {
      this.dialogRef.updateSize('400px');
    }
  }

  ngOnDestroy() {
    this.form.get('password').setValue('');
    this.form.get('confirm_password').setValue('');

    this.passwordSubmit.complete();

    this.subscriptions.forEach(sub => sub.unsubscribe());
  }

  proceed() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();
    this.disableDismiss = true;

    this.passwordSubmit.next({
      password: this.form.get('password').value,
      close: this.close.bind(this),
      error: this.error.bind(this),
    });
  }

  private validateForm() {
    if (this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value.length === 0) {
        return { Required: true };
      }

      if (this.data.confirm && this.form.get('password').value !== this.form.get('confirm_password').value) {
        return { NotEqual: true };
      }
    }

    return null;
  }

  private close() {
    this.dialogRef.close();
  }

  private error(error: any) {
    if (typeof error === 'object') {
      switch (error.status) {
        case 400:
          error = parseResponseMessage(error['_body']);
          break;
        case 401:
          error = 'Incorrect password';
          break;
        case 403:
          error = 'API Disabled';
          break;
        case 404:
          error = 'Wallet does not exist';
          break;
        default:
          const config = new MatSnackBarConfig();
          config.duration = 5000;
          this.snackbar.open(parseResponseMessage(error['_body']), null, config);
      }
    }

    this.button.setError(error ? error : 'Incorrect password');
    this.disableDismiss = false;
  }
}
