import { test, expect } from '@playwright/test';

test.describe('File Download Operations', () => {
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
