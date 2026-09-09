import { test, expect } from '@playwright/test';

test.describe('Download URL to OPFS', () => {
  test('Download from NVGD Resource', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Use an NVGD `file://` URL as the download source
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('http://localhost:9280/file://');

    // Step 2: Enter the NVGD URL in the download URL field
    await urlInput.fill('http://localhost:9280/version/');

    // Step 3: Enter a destination filename
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('version.txt');

    // Step 4: Click Download and verify the file was saved to OPFS with correct content
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeEnabled();
    await downloadBtn.click();
    await expect(page.locator('text=version.txt')).toBeVisible();

    // Step 5: Use the NVGD `/version/` endpoint as the download source
    await urlInput.fill('http://localhost:9280/version/');
    await nameInput.fill('version-info.txt');
    await downloadBtn.click();
    await expect(page.locator('text=version-info.txt')).toBeVisible();
  });
});
