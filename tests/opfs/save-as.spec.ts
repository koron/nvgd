// spec: specs/opfs-ui-test-plan.md
// TC-23: Save as — ファイルのローカル保存（正常系）
//
// window.showSaveFilePicker は Chromium のみのネイティブ API であり、
// Firefox / WebKit では未定義のため、addInitScript でモックに差し替えて
// 全ブラウザで実行できるようにする。

import { test, expect } from '@playwright/test';
import { createOPFSFile, reloadListing } from './helpers';

/** showSaveFilePicker のモックを page に注入する */
async function injectSaveFilePickerMock(page: import('@playwright/test').Page): Promise<void> {
  await page.addInitScript(() => {
    (window as any).__savedFiles = {} as Record<string, string>;
    (window as any).showSaveFilePicker = async ({
      suggestedName,
    }: {
      suggestedName: string;
    }) => {
      const chunks: BlobPart[] = [];
      return {
        name: suggestedName,
        createWritable: async () => ({
          write: async (data: BlobPart) => {
            chunks.push(data);
          },
          close: async () => {
            const blob = new Blob(chunks);
            (window as any).__savedFiles[suggestedName] = await blob.text();
          },
        }),
      };
    };
  });
}

test.describe('TC-23: Save as — ファイルのローカル保存（正常系）', () => {
  test('Save as リンクをクリックするとファイル内容が保存され成功アラートが表示される', async ({
    page,
  }) => {
    // モックを注入してからページに移動する（addInitScript は goto の前に呼ぶ）
    await injectSaveFilePickerMock(page);
    await page.goto('/opfs/');

    // テスト用ファイルを OPFS に作成してリロード
    await createOPFSFile(page, 'save-test.txt', 'Save as test content');
    await reloadListing(page);

    // Save as リンクが表示されている
    const row = page.locator('.grid-row').filter({ hasText: 'save-test.txt' });
    await expect(row.getByText('Save as')).toBeVisible();

    // Save as をクリック → 成功アラートが表示される
    const dialogPromise = page.waitForEvent('dialog');
    await row.getByText('Save as').click();
    const dialog = await dialogPromise;
    expect(dialog.message()).toContain('save-test.txt');
    expect(dialog.message()).toContain('successfully');
    await dialog.accept();

    // モックの __savedFiles にファイル内容が書き込まれている
    const saved = await page.evaluate(
      () => (window as any).__savedFiles['save-test.txt'] as string,
    );
    expect(saved).toBe('Save as test content');
  });
});
