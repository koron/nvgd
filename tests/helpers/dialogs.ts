import type { Dialog, Page } from '@playwright/test';

/**
 * Wait for the next dialog (alert/confirm/prompt) and respond to it.
 * Returns the captured dialog so the caller can assert on .message().
 *
 * Usage:
 *   const dialogPromise = captureDialog(page, 'accept');
 *   await page.click('#command-delete');
 *   const dialog = await dialogPromise;
 *   expect(dialog.message()).toMatch(/Are you sure/);
 */
export function captureDialog(
  page: Page,
  action: 'accept' | 'dismiss' = 'accept',
): Promise<Dialog> {
  return new Promise<Dialog>((resolve) => {
    page.once('dialog', async (d) => {
      try {
        if (action === 'accept') {
          await d.accept();
        } else {
          await d.dismiss();
        }
      } finally {
        resolve(d);
      }
    });
  });
}

/**
 * Auto-accept every dialog for the lifetime of the returned disposer.
 * Useful when a workflow chains multiple dialogs (e.g. confirm followed
 * by an informational alert).
 */
export function autoAcceptDialogs(page: Page): () => void {
  const handler = async (d: Dialog) => {
    await d.accept();
  };
  page.on('dialog', handler);
  return () => page.off('dialog', handler);
}

/**
 * Collect every dialog message emitted during a block. The returned
 * cleanup function detaches the handler and returns the accumulated
 * messages in order.
 */
export function recordDialogs(
  page: Page,
  action: 'accept' | 'dismiss' = 'accept',
): { messages: string[]; stop: () => string[] } {
  const messages: string[] = [];
  const handler = async (d: Dialog) => {
    messages.push(d.message());
    try {
      if (action === 'accept') {
        await d.accept();
      } else {
        await d.dismiss();
      }
    } catch { /* dialog may already be handled */ }
  };
  page.on('dialog', handler);
  return {
    messages,
    stop: () => {
      page.off('dialog', handler);
      return messages;
    },
  };
}
