import { mkdirSync, writeFileSync, existsSync } from 'node:fs';
import { join, dirname } from 'node:path';

/**
 * Generate fixture artifacts that we don't want to commit to the repo.
 *
 * - `large.bin`: 64KiB + 1 byte file. Used to verify the "Edit" action
 *   is hidden for files >= 64KiB (see D5).
 */
async function globalSetup(): Promise<void> {
  const fixturesDir = join(__dirname, 'fixtures');
  mkdirSync(fixturesDir, { recursive: true });

  const largePath = join(fixturesDir, 'large.bin');
  if (!existsSync(largePath)) {
    const buf = Buffer.alloc(64 * 1024 + 1, 0);
    writeFileSync(largePath, buf);
  }
}

export default globalSetup;
