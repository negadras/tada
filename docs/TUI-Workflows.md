# TUI User Workflows & Examples

## Overview

This guide provides practical examples of how to use the Tada TUI effectively in real-world scenarios. It covers common workflows, tips, and best practices for productivity.

## Table of Contents

- [Getting Started](#getting-started)
- [Daily Workflows](#daily-workflows)
- [Advanced Usage](#advanced-usage)
- [Tips & Tricks](#tips--tricks)
- [Troubleshooting](#troubleshooting)

## Getting Started

### First Time Setup

1. **Launch the TUI**
   ```bash
   tada --tui
   ```

2. **Explore the Dashboard**
   - Review statistics
   - Navigate menu items with ↑/↓
   - Press Enter to select

3. **Add Your First Todo**
   - Navigate to "Todo Management" → Enter
   - Press `a` to add
   - Fill in description and priority
   - Press Enter to save

4. **Add Your First Quote**
   - Navigate to "Quote Collection" → Enter
   - Press `a` to add
   - Fill in quote text, author, and category
   - Press Enter to save

## Daily Workflows

### 🌅 **Morning Routine - Planning Your Day**

#### Scenario: Start your day by reviewing and organizing tasks

```bash
# Launch TUI at dashboard
tada --tui
```

**Step-by-step:**
1. **Review Dashboard Statistics**
   - Check completion rate from yesterday
   - Note today's completed tasks
   - Assess overall productivity

2. **Navigate to Todos** (Tab or Enter on Todo Management)
   - Press `f` to filter by status
   - Review open tasks
   - Prioritize high-priority items

3. **Add Today's Tasks**
   - Press `a` to add new todo
   - Enter task description
   - Set appropriate priority
   - Repeat for multiple tasks

4. **Organize Existing Tasks**
   - Select tasks with ↑/↓
   - Press `e` to edit priorities
   - Update descriptions as needed

**Key Tips:**
- Use `f` to filter by priority: see high-priority tasks first
- Press `t` or `Enter` to quickly toggle task completion
- Add context to descriptions (e.g., "Call dentist @phone")

### 🌇 **Evening Routine - Reviewing Progress**

#### Scenario: End-of-day review and planning for tomorrow

```bash
# Launch directly to todos
tada --tui --screen todos
```

**Step-by-step:**
1. **Review Completed Tasks**
   - Press `f` to filter by "Done"
   - Review what you accomplished
   - Feel good about your progress!

2. **Update Task Status**
   - Press `f` to show "All" tasks
   - Mark completed tasks as done with `t`
   - Edit tasks that need updates

3. **Plan for Tomorrow**
   - Add new tasks for tomorrow
   - Set priorities based on deadlines
   - Add notes to existing tasks

4. **Get Inspiration**
   - Navigate to Quotes (Tab)
   - Press `Space` for random motivation
   - Add new quotes you found today

### 📊 **Weekly Review - Analyzing Productivity**

#### Scenario: Weekly productivity analysis and planning

```bash
# Start with dashboard overview
tada --tui --screen dashboard
```

**Step-by-step:**
1. **Analyze Statistics**
   - Review completion rates
   - Note productivity patterns
   - Identify areas for improvement

2. **Review All Tasks**
   - Navigate to Todos
   - Press `f` to see all tasks
   - Look for patterns in task types

3. **Clean Up Completed Tasks**
   - Filter by "Done" status
   - Review old completed tasks
   - Delete outdated tasks with `d`

4. **Plan Next Week**
   - Add recurring tasks
   - Set weekly goals
   - Adjust priorities based on deadlines

## Advanced Usage

### 🎯 **Project Management**

#### Scenario: Managing a multi-task project

**Organization Strategy:**
- Use descriptive task names: "Project X: Setup database"
- Use priority levels to indicate urgency
- Add context with tags in descriptions

**Workflow:**
1. **Create Project Tasks**
   ```
   Add multiple related tasks:
   - "Project X: Research requirements" (High)
   - "Project X: Design architecture" (Medium)
   - "Project X: Implement core features" (High)
   - "Project X: Write tests" (Medium)
   - "Project X: Documentation" (Low)
   ```

2. **Track Progress**
   - Use dashboard to monitor completion rate
   - Filter by priority to focus on urgent tasks
   - Update task descriptions with progress notes

3. **Manage Dependencies**
   - Edit task descriptions to note dependencies
   - Use priority levels to sequence tasks
   - Mark blocked tasks with lower priority

### 📚 **Knowledge Management with Quotes**

#### Scenario: Building a personal knowledge base

**Organization Strategy:**
- Use categories for different topics
- Include author information for credibility
- Add quotes as you learn new concepts

**Workflow:**
1. **Categorize Knowledge**
   ```
   Categories to use:
   - "productivity" - Work efficiency tips
   - "leadership" - Management insights
   - "learning" - Educational quotes
   - "motivation" - Inspirational content
   - "technical" - Programming wisdom
   ```

2. **Regular Addition**
   - Add quotes as you read articles/books
   - Include author and source when possible
   - Use consistent categorization

3. **Review and Reflect**
   - Use `f` to filter by category
   - Press `Space` for random inspiration
   - Press `Enter` to read full quotes

### 🔄 **Context Switching Between CLI and TUI**

#### Scenario: Seamless workflow between command-line and interface

**Quick Tasks (CLI):**
```bash
# Add tasks quickly
tada add "Call client about project status" --priority high
tada add "Review pull request #123" --priority medium

# Check specific information
tada list --status done --priority high
```

**Complex Tasks (TUI):**
```bash
# Interactive editing
tada update 5 --tui         # Edit task #5 interactively
tada quote update 3 --tui   # Edit quote #3 interactively

# Bulk management
tada --tui --screen todos   # Manage multiple todos
```

**Integration Tips:**
- Use CLI for automation and scripts
- Use TUI for exploration and bulk operations
- Switch contexts based on task complexity

## Tips & Tricks

### ⌨️ **Keyboard Efficiency**

#### Navigation Shortcuts
- `Tab` / `Shift+Tab` - Navigate between screens quickly
- `↑` / `↓` - Navigate lists efficiently
- `Enter` - Select items or submit forms
- `Esc` - Cancel actions or go back
- `?` - Show help when confused

#### Task Management Shortcuts
- `a` - Add new item (works in any list)
- `e` - Edit selected item
- `d` - Delete with confirmation
- `t` - Toggle todo status quickly
- `f` - Filter items by status/category
- `Space` - Special actions (random quote, etc.)

### 🎨 **Visual Organization**

#### Priority System
- **High Priority**: Critical tasks, deadlines
- **Medium Priority**: Important but not urgent
- **Low Priority**: Nice-to-have, future tasks

#### Description Patterns
```
Good task descriptions:
✓ "Fix login bug in user authentication"
✓ "Meeting with team @2pm about project X"
✓ "Research React hooks for dashboard component"

Poor task descriptions:
✗ "Fix bug"
✗ "Meeting"
✗ "Research"
```

#### Category Organization
```
Quote categories:
- "work" - Professional development
- "life" - Personal growth
- "tech" - Technical insights
- "leadership" - Management wisdom
- "creativity" - Inspiration for creative work
```

### 📱 **Form Efficiency**

#### Quick Form Navigation
- `Tab` - Move to next field
- `Shift+Tab` - Move to previous field
- `Enter` - Submit form
- `Esc` - Cancel form

#### Form Best Practices
- Fill required fields first (marked with *)
- Use consistent naming conventions
- Keep descriptions concise but informative
- Use categories consistently

### 🔍 **Search and Filter Strategies**

#### Todo Filtering
- Filter by status to focus on open tasks
- Filter by priority during busy periods
- Use "All" filter for comprehensive review

#### Quote Filtering
- Filter by category for specific inspiration
- Use author names to find specific quotes
- Browse all quotes for serendipitous discovery

## Troubleshooting

### 🐛 **Common Issues**

#### Screen Display Problems
**Issue:** TUI doesn't display properly
**Solution:**
- Check terminal size (minimum 80x24)
- Ensure terminal supports colors
- Try different terminal emulators

#### Database Connection Issues
**Issue:** "Failed to open database" error
**Solution:**
- Check file permissions
- Ensure database directory exists
- Restart the application

#### Form Submission Issues
**Issue:** Form doesn't submit
**Solution:**
- Ensure all required fields are filled
- Check for validation errors
- Use `Enter` to submit, not `Esc`

### 🛠️ **Performance Tips**

#### Large Datasets
- Use filters to reduce displayed items
- Delete old completed tasks periodically
- Consider archiving old quotes

#### Slow Startup
- Check database file size
- Ensure adequate system resources
- Close other resource-intensive applications

#### Terminal Settings
- Enable 256 colors
- Set font to monospace
- Configure appropriate size (80x24 minimum)

## Best Practices

### 📋 **Task Management**

1. **Regular Review**
   - Daily: Review and update tasks
   - Weekly: Clean up completed tasks
   - Monthly: Analyze productivity patterns

2. **Consistent Naming**
   - Use action verbs in task descriptions
   - Include context when helpful
   - Be specific about deliverables

3. **Priority Management**
   - Use high priority sparingly
   - Adjust priorities based on deadlines
   - Consider both urgency and importance

### 💬 **Quote Collection**

1. **Quality over Quantity**
   - Add quotes that resonate with you
   - Include proper attribution
   - Use meaningful categories

2. **Regular Engagement**
   - Use random quote feature for inspiration
   - Review specific categories when needed
   - Add new quotes as you discover them

3. **Organization**
   - Use consistent category names
   - Include author information
   - Remove duplicates periodically

### 🔄 **Workflow Integration**

1. **Choose the Right Tool**
   - CLI for quick additions
   - TUI for complex management
   - Both for different contexts

2. **Maintain Consistency**
   - Use same categories across CLI and TUI
   - Keep naming conventions consistent
   - Sync workflows between interfaces

3. **Automation Where Possible**
   - Use CLI in scripts
   - Integrate with other tools
   - Create aliases for common commands
