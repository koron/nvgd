// spec: specs/opfs-ui-test-plan.md
// TC-20: DuckDB 連携 — 対応形式ファイル
// TC-21: DuckDB 連携 — 非対応形式ファイル

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, reloadListing } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-20: DuckDB 連携 — 対応形式ファイル', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('CSV ファイルを選択して DuckDB シェルを開ける', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'data.csv', 'id,name\n1,Alice\n');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="data.csv"]').check();
    await expect(page.locator('#command-duckdb')).toBeEnabled();

    // 新しいタブ（popup）が開く
    const [popup] = await Promise.all([
      page.waitForEvent('popup'),
      page.click('#command-duckdb'),
    ]);

    // /duckdb/ に遷移して opfs= パラメーターが含まれる
    await expect(popup).toHaveURL(/\/duckdb\//);
    expect(popup.url()).toContain('opfs=');
  });
});

test.describe('TC-21: DuckDB 連携 — 非対応形式ファイル', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('txt ファイルを選択すると DuckDB シェルは開くが VIEW が作られない URL になる', async ({
    page,
  }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'note.txt', 'plain text');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="note.txt"]').check();
    await expect(page.locator('#command-duckdb')).toBeEnabled();

    const [popup] = await Promise.all([
      page.waitForEvent('popup'),
      page.click('#command-duckdb'),
    ]);

    // /duckdb/ には遷移する
    await expect(popup).toHaveURL(/\/duckdb\//);

    // opfs= パラメーターは含まれるが、ハッシュ部分に CREATE VIEW がない
    const url = popup.url();
    expect(url).toContain('opfs=');
    expect(url).not.toContain('CREATE+VIEW');
    expect(url).not.toContain('CREATE%20VIEW');
  });
});
