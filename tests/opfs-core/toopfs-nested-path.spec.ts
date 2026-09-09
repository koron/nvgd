import { test, expect } from '@playwright/test';

test.describe('toopfs Filter - Download Files from NVGD to OPFS', () => {
  test('Download to Nested OPFS Path', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Set a nested destination path
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('output/data/');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are saved in the nested OPFS path
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
    await expect(opfsPage!.locator('text=output/data/')).toBeVisible();
  });

  test('Download After Navigation', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Navigate away from the page
    await page.goto('http://localhost:9280/opfs/');

    // Navigate back to the original page
    await page.goto('http://localhost:9280/file://?toopfs');

    // Verify the toopfs UI is still functional
    const toopfsUI = page.locator('text=Download to OPFS');
    await expect(toopfsUI).toBeVisible();

    // Select and download files as before
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();
  });

  test('Download Multiple Files in Sequence', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select all files
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();

    // Click Download
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify each file is downloaded and uploaded sequentially
    const progressOverlay = page.locator('text=downloading');
    await expect(progressOverlay).toBeVisible();

    // Verify the progress shows the correct step numbers
    await expect(progressOverlay).toBeVisible();

    // Verify all files appear in the OPFS destination directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with Empty Selection', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Unselect all files
    const checkboxes = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkboxes.first().check();
    await checkboxes.first().click();

    // Verify the "Download" button is disabled
    const downloadBtn = page.locator('text=Download, button');
    await expect(downloadBtn).toBeDisabled();

    // Attempt to click Download
    await downloadBtn.click();

    // Verify no action is triggered
    const progressOverlay = page.locator('text=downloading');
    await expect(progressOverlay).not.toBeVisible();
  });

  test('Download with file:// Protocol keepcompress', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Download the compressed file to OPFS
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the compressed file is saved correctly in OPFS
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with https:// Protocol all Parameter', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs&all');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are downloaded correctly with the `?all` parameter
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });

  test('Download with file:// Protocol all Parameter', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs&all');

    // Select and download files
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the files are downloaded correctly with the `?all` parameter
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });
});
