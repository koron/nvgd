import { test, expect } from '@playwright/test';

test.describe('Edge Cases and Error Handling', () => {
  test('OPFS Operations with Timeout', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Navigate to an NVGD page with a slow source URL and `?toopfs` filter
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('http://localhost:9280/help/');
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('timeout-test.md');

    // Step 2: Select a file and attempt to download
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();

    // Step 3: Verify the download completes or times out
    const errorMsg = page.locator('text=Error, text=Failed');
    await expect(errorMsg).toBeVisible();

    // Step 4: Verify the progress overlay is removed after timeout
    await expect(page.locator('text=downloading')).not.toBeVisible();
  });

  test('OPFS Operations with Invalid Input', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Enter an invalid directory name in the "Directory name to create" input
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('');
    await page.locator('text=Create new directory, button').click();
    const dirError = page.locator('text=Need directory name');
    await expect(dirError).toBeVisible();

    // Step 2: Enter an invalid file name in the "File name to create or update" field
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Invalid file name');
    await page.locator('text=Create or update a file, button').click();
    const fileError = page.locator('text=Need file name');
    await expect(fileError).toBeVisible();

    // Step 3: Enter an invalid URL in the "URL to be downloaded" field
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('not-a-valid-url');
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('invalid.md');
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    const urlError = page.locator('text=Error, text=Failed');
    await expect(urlError).toBeVisible();
  });
});
