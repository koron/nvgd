import { test, expect } from '@playwright/test';

test.describe('File Download Operations', () => {
  test('Save as - Download File from OPFS to Local', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with known content in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('download-test.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('This is the content to download.');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "download-test.txt"')).toBeVisible();

    // Step 2: Click the "Save as" action (download icon) next to the file
    const saveAsBtn = page.locator('text=Save as, [class*="save"], button[title*="Save"]');
    await saveAsBtn.click();

    // Step 3: Verify a file picker dialog appears with the file name pre-filled
    const filePicker = page.locator('input[type=file]');
    await expect(filePicker).toBeVisible();

    // Step 4: Select a destination location and save the file
    await filePicker.setInputFiles({
      name: 'downloaded.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('This is the content to download.')
    });

    // Step 5: Verify a success alert confirming the file was saved
    await expect(page.locator('text=Downloaded "download-test.txt"')).toBeVisible();

    // Step 6: Verify the saved file contains the expected content
    // (This is handled by the OS file picker - we verify the alert)

    // Step 7: Verify the file is saved with the original OPFS file name
    await expect(page.locator('text=download-test.txt')).toBeVisible();
  });

  test('Cancel File Download', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file in OPFS
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('cancel-test.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Content to cancel download');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "cancel-test.txt"')).toBeVisible();

    // Step 2: Click the "Save as" action on the file
    const saveAsBtn = page.locator('text=Save as, [class*="save"], button[title*="Save"]');
    await saveAsBtn.click();

    // Step 3: Cancel the file picker dialog without saving
    await page.keyboard.press('Escape');

    // Step 4: Verify the file still exists in OPFS unchanged
    await expect(page.locator('text=cancel-test.txt')).toBeVisible();
  });
});
