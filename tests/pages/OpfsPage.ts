import { expect, type Locator, type Page } from '@playwright/test';
import { clearOPFS } from '../helpers/opfs';

/**
 * Page Object for the OPFS Web UI ( /opfs/ ).
 *
 * Centralises every locator and high-level action used by the spec
 * files so individual tests stay readable.
 */
export class OpfsPage {
  readonly page: Page;

  // Header / structure
  readonly title: Locator;
  readonly breadcrumb: Locator;
  readonly grid: Locator;
  readonly toggleSelectAll: Locator;

  // Footer command buttons
  readonly reloadBtn: Locator;
  readonly deleteBtn: Locator;
  readonly duckdbBtn: Locator;

  // mkdir
  readonly mkdirName: Locator;
  readonly mkdirBtn: Locator;

  // upload
  readonly uploadFileInput: Locator;
  readonly uploadName: Locator;
  readonly uploadBtn: Locator;

  // editor
  readonly editorName: Locator;
  readonly editorBody: Locator;
  readonly editorSaveBtn: Locator;
  readonly editorClearBtn: Locator;

  // URL download
  readonly downloadUrl: Locator;
  readonly downloadAs: Locator;
  readonly downloadBtn: Locator;
  readonly downloadClearBtn: Locator;

  constructor(page: Page) {
    this.page = page;

    this.title = page.locator('title');
    this.breadcrumb = page.locator('#header');
    this.grid = page.locator('#main > .directory');
    this.toggleSelectAll = page.locator('#toggle-selection-all');

    this.reloadBtn = page.locator('#command-reload');
    this.deleteBtn = page.locator('#command-delete');
    this.duckdbBtn = page.locator('#command-duckdb');

    this.mkdirName = page.locator('#mkdir-name');
    this.mkdirBtn = page.locator('#mkdir-mkdir');

    this.uploadFileInput = page.locator('#upload-file');
    this.uploadName = page.locator('#upload-name');
    this.uploadBtn = page.locator('#upload-upload');

    this.editorName = page.locator('#editor-name');
    this.editorBody = page.locator('#editor-edit');
    this.editorSaveBtn = page.locator('#editor-save');
    this.editorClearBtn = page.locator('#editor-clear');

    this.downloadUrl = page.locator('#download-url');
    this.downloadAs = page.locator('#download-as');
    this.downloadBtn = page.locator('#download-download');
    this.downloadClearBtn = page.locator('#download-clear');
  }

  /**
   * Navigate to /opfs/ (optionally with a hash path like `sub1/sub2/`).
   * Waits until the grid header is rendered.
   *
   * Note: when the previous URL is already /opfs/ and we only change
   * the hash, the browser treats this as a same-document navigation
   * and does NOT re-run the page-s init() function. We force a reload
   * after hashed navigations so init() picks up the new path.
   */
  async goto(hashPath: string = ''): Promise<void> {
    const url = hashPath ? `/opfs/#${hashPath}` : '/opfs/';
    await this.page.goto(url);
    if (hashPath) {
      await this.page.reload();
    }
    await expect(this.grid.locator('.grid-header')).toBeVisible();
  }

  /** Convenience: go to /opfs/ and wipe OPFS, then reload. */
  async gotoAndReset(): Promise<void> {
    await this.goto();
    await clearOPFS(this.page);
    await this.reloadBtn.click();
  }

  /** Read the current breadcrumb as an array of segment texts. */
  async breadcrumbSegments(): Promise<string[]> {
    return await this.breadcrumb.evaluate((el) => {
      return Array.from(el.children)
        .filter((c) => c.tagName === 'A' || c.tagName === 'SPAN')
        .map((c) => c.textContent ?? '')
        .filter((t) => t.trim() !== '/' && t.trim() !== '');
    });
  }

  /** Locator for a specific row by display name (e.g. `foo/` or `bar.txt`). */
  row(displayName: string): Locator {
    return this.grid.locator('.grid-row', {
      has: this.page.locator(`input.selectedFile[name="${displayName}"]`),
    });
  }

  /** All file/directory names visible in the current directory listing. */
  async rowNames(): Promise<string[]> {
    return await this.grid
      .locator('input.selectedFile')
      .evaluateAll((inputs) =>
        inputs.map((i) => (i as HTMLInputElement).name),
      );
  }

  /** Toggle the checkbox for the given row. */
  async selectRow(displayName: string): Promise<void> {
    await this.row(displayName).locator('input.selectedFile').check();
  }

  async unselectRow(displayName: string): Promise<void> {
    await this.row(displayName).locator('input.selectedFile').uncheck();
  }

  /** Click the directory link in a row (navigates into it). */
  async openDirectory(name: string): Promise<void> {
    await this.row(`${name}/`).locator('a').click();
    await expect(this.grid.locator('.grid-header')).toBeVisible();
  }

  /** Create a directory via the mkdir form. */
  async mkdir(name: string): Promise<void> {
    await this.mkdirName.fill(name);
    await this.mkdirBtn.click();
  }

  /**
   * Upload a local file. If `asName` is provided it overrides the
   * suggested filename in the upload form.
   */
  async uploadLocalFile(localPath: string, asName?: string): Promise<void> {
    await this.uploadFileInput.setInputFiles(localPath);
    if (asName !== undefined) {
      await this.uploadName.fill(asName);
    }
    await this.uploadBtn.click();
  }

  /** Create or update a file via the simple editor. */
  async saveEditor(name: string, body: string): Promise<void> {
    await this.editorName.fill(name);
    await this.editorBody.fill(body);
    await this.editorSaveBtn.click();
  }

  /**
   * Click the "Edit" action on a file row. The action is rendered as
   * `<a><span class="material-symbols">edit</span>Edit</a>`, so its
   * textContent is `editEdit` -- `getByText("Edit", exact: true)`
   * will not match. We pick the `<a>` whose text *contains* "Edit".
   */
  async clickEdit(fileName: string): Promise<void> {
    await this.row(fileName)
      .locator('a')
      .filter({ hasText: 'Edit' })
      .click();
  }

  /** Click the "Save as" action on a file row. */
  async clickSaveAs(fileName: string): Promise<void> {
    await this.row(fileName)
      .locator('a')
      .filter({ hasText: 'Save as' })
      .click();
  }

  /** Locator helpers for asserting on action visibility. */
  editAction(fileName: string): Locator {
    return this.row(fileName).locator('a').filter({ hasText: 'Edit' });
  }

  saveAsAction(fileName: string): Locator {
    return this.row(fileName).locator('a').filter({ hasText: 'Save as' });
  }

  /** Fill in the URL download form and submit. */
  async downloadFromUrl(url: string, asName: string): Promise<void> {
    await this.downloadUrl.fill(url);
    await this.downloadAs.fill(asName);
    await this.downloadBtn.click();
  }

  /** Number of file/directory rows currently rendered. */
  async rowCount(): Promise<number> {
    return await this.grid.locator('.grid-row').count();
  }
}
