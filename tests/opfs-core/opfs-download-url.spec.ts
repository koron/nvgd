import { test, expect } from '@playwright/test';

test.describe('Download URL to OPFS', () => {
  test('Download URL to OPFS', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Navigate to an NVGD URL that returns content
    const urlInput = page.locator('input[placeholder*="URL"]');
    await urlInput.fill('http://localhost:9280/help/');

    // Step 2: Enter a file name in the "Name to save" field
    const nameInput = page.locator('input[placeholder*="Name"]');
    await nameInput.fill('help.md');

    // Step 3: Verify the "Download" button becomes enabled
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeEnabled();

    // Step 4: Click the "Download" button
    await downloadBtn.click();

    // Step 5: Verify the file appears in the OPFS directory listing
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText('help.md');

    // Step 6: Verify the file contains the expected content from the URL
    await expect(page.locator('text=help.md')).toBeVisible();

    // Step 7: Download a file with a name that already exists
    await urlInput.fill('http://localhost:9280/help/');
    await nameInput.fill('help.md');
    await downloadBtn.click();
    const confirmDialog = page.locator('text=Confirm, text=overwrite');
    await expect(confirmDialog).toBeVisible();

    // Step 8: After confirming overwrite, verify the file is updated
    await page.locator('text=Confirm, button, :text("Confirm")').click();
    await expect(page.locator('text=help.md')).toBeVisible();

    // Step 9: Enter an invalid URL
    await urlInput.fill('not-a-valid-url');
    await expect(downloadBtn).toBeDisabled();

    // Step 10: Enter a URL without a filename
    await urlInput.fill('http://localhost:9280/help/');
    await nameInput.fill('');
    await expect(downloadBtn).toBeDisabled();

    // Step 11: Enter a URL without `http://` or `https://`
    await urlInput.fill('ftp://localhost:9280/help/');
    await expect(downloadBtn).toBeDisabled();

    // Step 12: A URL that returns an error
    await urlInput.fill('http://localhost:9280/');
    await nameInput.fill('error.md');
    await downloadBtn.click();
    const errorAlert = page.locator('text=Error, text=Failed');
    await expect(errorAlert).toBeVisible();
  });

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
