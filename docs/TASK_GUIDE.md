# Task Guide

## Overview
Tasks track units of work from start to completion. Each task gets a markdown file.

## Directory Structure
```
task/
  active/      -- Tasks currently being worked on
  done/        -- Completed tasks
  archive/     -- Old tasks for reference
  templates/   -- Task templates
    GENERIC.md
    golang/
      DESIGN.md
      IMPLEMENTATION.md
```

## Workflow

### 1. Create Task
```bash
cp task/templates/golang/IMPLEMENTATION.md task/active/<task-name>.md
# Edit the file to fill in details
```

### 2. Work on Task
- Update checklist items as you progress
- Note blockers or decisions in the task file

### 3. Complete Task
```bash
mv task/active/<task-name>.md task/done/
```

### 4. Archive (periodically)
```bash
mv task/done/<old-task>.md task/archive/
```

## Task File Contents
Every task file should include:
- **Summary**: What and why
- **Scope**: Files/packages affected
- **Acceptance Criteria**: How to know it is done
- **Checklist**: Step-by-step items
- **Notes**: Decisions, blockers, learnings

## Naming Convention
`<YYYY-MM-DD>-<short-description>.md`

Example: `2026-03-06-add-batch-fetch-command.md`
