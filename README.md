# d-day-labeler

**d-day-labeler** is a GitHub Action that automatically decrements 'D-n' style labels on pull requests every night at midnight. It helps teams manage deadlines or remaining days easily by updating the labels.

## Features

- **Auto Decrement Labels**: Every night at midnight (KST), 'D-n' style labels on open pull requests are decremented by 1 (e.g., D-1 → D-0).
- **Easy Deadline Tracking**: Helps keep track of deadlines by automatically adjusting labels such as D-7, D-6, etc.
- **Minimal Setup**: Simple GitHub Action integration with minimal configuration required.

## Usage

Create a `.github/workflows/decrement-labels.yml` file in your repository with the following content:

```yaml
name: Update Day Labels
on:
  schedule:
    - cron: '0 0 * * *' # Runs every midnight KST.
jobs:
  decrement-labels:
    runs-on: ubuntu-latest
    steps:
      - name: Update D-n Labels
        uses: devmyong/d-day-labeler@v1.0.1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

- `token` (required):
  - **Description**: GitHub API token to authenticate the action.
  - **Example**: `${{ secrets.GITHUB_TOKEN }}`

## Example Labels

- D-7, D-6, D-5, ..., D-1, D-0

### License

[MIT License](LICENSE)
