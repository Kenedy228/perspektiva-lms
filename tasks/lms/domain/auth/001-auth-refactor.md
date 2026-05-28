# User Roles and Permissions

# Overview

The LMS system contains four roles:

* Administrator
* Creator
* Organization
* Student

The system uses strict role-based access control (RBAC).

All permissions must be validated on the backend.

---

# Permission Matrix

| Action                           | Administrator | Creator | Organization                 | Student                    |
| -------------------------------- | ------------- | ------- | ---------------------------- | -------------------------- |
| Sign in                          | yes           | yes     | yes                          | yes                        |
| Sign out                         | yes           | yes     | yes                          | yes                        |
| Create user account              | yes           | no      | no                           | no                         |
| View user accounts               | yes           | no      | no                           | no                         |
| Edit/Delete user account         | yes           | no      | no                           | no                         |
| Create course                    | yes           | yes     | no                           | no                         |
| View courses                     | yes           | yes     | no                           | yes (assigned only)        |
| Edit/Delete course               | yes           | yes     | no                           | no                         |
| Add block/element to course      | yes           | yes     | no                           | no                         |
| View course block/element        | yes           | yes     | no                           | yes (assigned course only) |
| Edit/Delete course block/element | yes           | yes     | no                           | no                         |
| Create question bank             | yes           | yes     | no                           | no                         |
| View question bank               | yes           | yes     | no                           | no                         |
| Edit/Delete question bank        | yes           | yes     | no                           | no                         |
| Create question                  | yes           | yes     | no                           | no                         |
| View questions                   | yes           | yes     | no                           | no                         |
| Edit/Delete question             | yes           | yes     | no                           | no                         |
| Create test                      | yes           | yes     | no                           | no                         |
| View test                        | yes           | yes     | no                           | yes (assigned course only) |
| Pass test                        | no            | no      | no                           | yes                        |
| Edit/Delete test                 | yes           | yes     | no                           | no                         |
| View student progress            | yes           | no      | yes (assigned students only) | yes (own only)             |

---

# Role Definitions

## Administrator

The administrator has unrestricted access to the entire system.

### Responsibilities

* Manage users
* Assign courses to students
* Manage organizations
* Manage educational content
* Monitor student progress
* Access all analytics

### Permissions Scope

Global system access.

---

## Creator

The creator manages educational content only.

### Responsibilities

* Create and manage courses
* Create and manage blocks
* Create and manage course elements
* Create and manage question banks
* Create and manage questions
* Create and manage tests

### Restrictions

* Cannot manage users
* Cannot assign students to courses
* Cannot access analytics
* Cannot view student progress

### Permissions Scope

Educational content management only.

---

## Organization

The organization role is analytical and read-only.

### Responsibilities

* View assigned students progress
* Monitor educational results
* Access student analytics

### Restrictions

* Cannot modify educational content
* Cannot manage users
* Cannot assign courses
* Cannot modify tests or attempts

### Permissions Scope

Only assigned students attached by administrator.

---

## Student

The student role is intended for learning process participation.

### Responsibilities

* View assigned courses
* Open course content
* Pass tests
* View own progress

### Restrictions

* Cannot modify educational content
* Cannot manage users
* Cannot access analytics of other students

### Permissions Scope

Only courses explicitly assigned by administrator.

---

# Access Rules

## Course Access

Students can access only courses assigned by administrator.

Unassigned courses must be completely hidden from students.

---

## Progress Visibility

* Administrators can view all student progress.
* Organizations can view only assigned students.
* Students can view only their own progress.
* Creators cannot access progress data.

---

# Backend Authorization Rules

## Important

Frontend authorization checks are NOT considered secure.

Every protected backend endpoint must validate:

1. Authentication
2. User role
3. Resource access scope

---

# Architecture Constraints

## Forbidden

* Frontend-only authorization
* Global permission bypasses
* Hidden role assumptions
* Authorization inside infrastructure layer

## Required

* Explicit permission checks
* Role validation inside application/use-case layer
* Clear separation between learning process and content management
* Resource ownership validation

---

# Notes For AI Agents

When implementing new features:

* Never expand permissions implicitly.
* Always follow least privilege principle.
* Students must never access unassigned resources.
* Organization role is strictly analytical.
* Creator role is strictly content-oriented.
* Administrator is the only role with full system access.
