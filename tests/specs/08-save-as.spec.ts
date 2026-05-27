import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog } from '../helpers/dialogs';
import { clearOPFS, seedFile } from '../helpers/opfs';

/**
 * H. Save as (Chromium only).
 *
 * `window.showSaveFilePicker` is only available in Chromium-based
 * browsers (Firefox has no File System Access API for write). We stub
 * the picker via addInitScript so the test never blocks on a real
 * native dialog.
 */

test.describe('H. Save as (Chromium only)', () => {
  test.skip(
    ({ browserName }) => browserName !== 'chromium',
    'showSaveFilePicker is Chromium-only',
  );

  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);

    // Install a mock for showSaveFilePicker so we never hit the OS
    // native dialog. The mock records writes on window.__writes.
    await page.addInitScript(() => {
      (window as any).__writes = [] as Uint8Array[];
      (window as any).__abortPicker = false;
      (window as any).showSaveFilePicker = async () => {
        if ((window as any).__abortPicker) {
          const err = new DOMException(
            'user aborted',
            'AbortError',
          );
          throw err;
        }
        return {
          name: 'mocked-output.txt',
          createWritable: async () => ({
            write: async (blob: Blob) => {
              const buf = new Uint8Array(await blob.arrayBuffer());
              (window as any).__writes.push(buf);
            },
            close: async () => {},
          }),
        };
      };
    });
  });

  test('H1: Save as writes the file through the picker mock', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'export.txt', 'export contents');
    await opfs.goto();

    // The success path ends with an info alert.
    const dialogPromise = captureDialog(page, 'accept');
    await opfs.clickSaveAs('export.txt');
    await dialogPromise;

    const writes = await page.evaluate(
      () => ((window as any).__writes as Uint8Array[]).map((u) => u.length),
    );
    expect(writes.length).toBe(1);
    expect(writes[0]).toBe('export contents'.length);
  });

  test('H2: cancelling the picker (AbortError) shows no alert', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'export.txt', 'x');
    await opfs.goto();

    await page.evaluate(() => {
      (window as any).__abortPicker = true;
    });

    let dialogFired = false;
    page.once('dialog', async (d) => {
      dialogFired = true;
      await d.accept();
    });

    await opfs.clickSaveAs('export.txt');
    // Give any potential alert a moment to surface.
    await page.waitForTimeout(300);

    expect(dialogFired).toBe(false);
  });
});
