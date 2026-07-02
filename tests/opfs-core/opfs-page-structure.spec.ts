import { test, expect } from '@playwright/test';

test.describe('OPFS Page Loading and Initial State', () => {
  test('OPFS Page Structure and Elements', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Verify the page contains a header section with breadcrumb trail
    const breadcrumb = page.locator('text=(Root)');
    await expect(breadcrumb).toBeVisible();

    // Verify the directory listing area (grid table)
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toBeVisible();

    // Verify footer controls: Reload, Delete, DuckDB buttons
    // Reload button
    await expect(page.locator('text=Reload, button, :text("Reload")')).toBeVisible();
    // Delete button (initially disabled)
    await expect(page.locator('text=Delete, button, :text("Delete")')).toBeVisible();
    // DuckDB button (initially disabled)
    await expect(page.locator('text=DuckDB, button, :text("DuckDB")')).toBeVisible();

    // Verify Make directory input and button
    const makeDirInput = page.locator('input[type=text], input[name*="dir"], input[placeholder*="dir"]');
    await expect(makeDirInput).toBeVisible();
    const makeDirButton = page.locator('text=Create new directory, button, :text("Create new directory")');
    await expect(makeDirButton).toBeVisible();

    // Verify Upload file input and button
    const uploadInput = page.locator('input[type=file]');
    await expect(uploadInput).toBeVisible();
    const uploadButton = page.locator('text=Upload, button, :text("Upload")');
    await expect(uploadButton).toBeVisible();

    // Verify Simple Editor section (name, textarea, save, clear buttons)
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await expect(editorName).toBeVisible();
    const editorTextarea = page.locator('textarea');
    await expect(editorTextarea).toBeVisible();
    const saveButton = page.locator('text=Create or update a file, button, :text("Create or update a file")');
    await expect(saveButton).toBeVisible();
    const clearButton = page.locator('text=Clear, button, :text("Clear")');
    await expect(clearButton).toBeVisible();

    // Verify Download URL section (URL input, name input, download button)
    const urlInput = page.locator('input[placeholder*="URL"]');
    await expect(urlInput).toBeVisible();
    const nameInput = page.locator('input[placeholder*="Name"]');
    await expect(nameInput).toBeVisible();
    const downloadButton = page.locator('text=Download, button, :text("Download")');
    await expect(downloadButton).toBeVisible();

    // Verify Delete and DuckDB buttons are initially disabled
    const deleteBtn = page.locator('text=Delete, button');
    await expect(deleteBtn).toBeDisabled();
    const duckdbBtn = page.locator('text=DuckDB, button');
    await expect(duckdbBtn).toBeDisabled();
    const dlBtn = page.locator('text=Download, button');
    await expect(dlBtn).toBeDisabled();
  });
});
