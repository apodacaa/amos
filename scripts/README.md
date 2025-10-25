# Performance Testing Scripts

## Test Data Generator

Generate synthetic entries and todos to test Amos performance at scale.

### Usage

```bash
# Light usage (100 entries, 50 todos)
go run scripts/generate_test_data.go -entries 100 -todos 50

# Heavy journal user (1,000 entries, 500 todos)
go run scripts/generate_test_data.go -entries 1000 -todos 500

# Stress test (10,000 entries, 5,000 todos)
go run scripts/generate_test_data.go -entries 10000 -todos 5000

# Test in separate directory (won't overwrite your data)
go run scripts/generate_test_data.go -entries 1000 -path /tmp/amos-test
```

### Flags

- `-entries N` - Number of entries to generate (default: 100)
- `-todos N` - Number of todos to generate (default: 50)
- `-path PATH` - Custom output path (default: ~/.amos)

### Generated Data

**Entries**:
- Realistic titles and bodies (200-800 characters)
- Random @tags (2-5 tags per entry)
- Timestamps spread over past year

**Todos**:
- 60% linked to random entries
- 40% standalone
- Random statuses (weighted towards "open")
- 1-3 tags per todo

### Testing Performance

1. **Backup your data** (if testing in ~/.amos):
   ```bash
   cp -r ~/.amos ~/.amos.backup
   ```

2. **Generate test data**:
   ```bash
   go run scripts/generate_test_data.go -entries 1000 -todos 500
   ```

3. **Run the app**:
   ```bash
   make run
   ```

4. **Test operations**:
   - Navigate entries list (j/k keys)
   - Filter by tag (@)
   - Toggle todo status (space)
   - Save new entry (Ctrl+S)
   - Observe lag/responsiveness

5. **Restore your data**:
   ```bash
   rm -rf ~/.amos
   mv ~/.amos.backup ~/.amos
   ```

### Expected Performance

Based on O(n) storage model:

| Entries | Todos | File Size | Load Time | Save Time | Status |
|---------|-------|-----------|-----------|-----------|--------|
| 100     | 50    | ~50 KB    | Instant   | Instant   | ✅ Excellent |
| 1,000   | 500   | ~500 KB   | <50ms     | ~100ms    | ✅ Good |
| 10,000  | 5,000 | ~5 MB     | ~200ms    | ~500ms    | ⚠️ Noticeable lag |
| 100,000 | 50,000| ~50 MB    | ~2s       | ~5s       | ❌ Unusable |

**Bottleneck**: Every save rewrites the entire JSON file (O(n) operation).

### Performance Observations

Document your findings:
- At what size do you notice lag?
- Which operations are slowest? (save vs load vs filter)
- File size at breaking point?
- Memory usage (check with `top` or Activity Monitor)?

This data informs whether optimizations are needed (caching, SQLite, etc.).
