// -*- Mode: Go; indent-tabs-mode: t -*-
// +build arm 386 ppc

/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package sys

import "syscall"

const (
	_SYS_GETUID  = syscall.SYS_GETUID32
	_SYS_GETGID  = syscall.SYS_GETGID32
	_SYS_GETEUID = syscall.SYS_GETEUID32
	_SYS_GETEGID = syscall.SYS_GETEGID32
)
