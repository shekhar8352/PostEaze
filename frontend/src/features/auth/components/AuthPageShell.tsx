import type { ReactNode } from 'react';
import styles from './AuthPageShell.module.css';

export type AuthPageVariant = 'login' | 'register' | 'forgot' | 'verify';

const copy: Record<
  AuthPageVariant,
  { kicker: string; headline: string; subline: string; meta: string }
> = {
  login: {
    kicker: 'Sign in',
    headline: 'Access your workspace.',
    subline:
      'Use your verified credentials. Session access is protected and audited for your organization.',
    meta: 'Enterprise-ready · Email verification required for new accounts.',
  },
  register: {
    kicker: 'Get started',
    headline: 'Create your organization account.',
    subline:
      'Verify your email once to activate access. You can invite teammates after onboarding.',
    meta: 'Data handling aligned with standard B2B practice · OAuth optional.',
  },
  forgot: {
    kicker: 'Recovery',
    headline: 'Reset your password.',
    subline:
      'We will send a signed link to your work email. Links expire automatically for security.',
    meta: 'If you do not receive mail, check spam or contact your IT administrator.',
  },
  verify: {
    kicker: 'Verification',
    headline: 'Confirm your email address.',
    subline:
      'This step protects your workspace and ensures only approved inboxes can activate accounts.',
    meta: 'You can resend the message or return to sign in once verification is complete.',
  },
};

type AuthPageShellProps = {
  variant: AuthPageVariant;
  children: ReactNode;
};

export function AuthPageShell({ variant, children }: AuthPageShellProps) {
  const c = copy[variant];

  return (
    <div className={styles.root}>
      <div className={styles.backdrop} aria-hidden />
      <div className={styles.grain} aria-hidden />
      <div className={styles.grid}>
        <aside className={styles.brand}>
          <div>
            <div className={styles.wordmark}>
              Post<span>Eaze</span>
            </div>
            <p className={styles.kicker}>{c.kicker}</p>
            <h1 className={styles.headline}>{c.headline}</h1>
            <p className={styles.subline}>{c.subline}</p>
            <div className={styles.stripe} aria-hidden />
          </div>
          <p className={styles.meta}>
            <strong>PostEaze</strong> — social publishing for teams that need clarity and control.
          </p>
        </aside>
        <main className={styles.main}>
          <div className={styles.card}>{children}</div>
        </main>
      </div>
    </div>
  );
}
