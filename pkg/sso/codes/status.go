/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package codes // import "go.pilab.hu/pilab-cloud/ga-pi/pkg/shadowsso/codes"

type Code uint32

const (
	// OK is returned on success.
	OK Code = iota

	// UserNotFound when account not found.
	UserNotFound

	// TooManyAttempts when too many attempts are made.
	TooManyAttempts

	// RequestMFA when 2FA is required.
	// Additional information should be provided in the response,
	// like the required method.
	RequestMFA

	// InvalidMFA when 2FA is invalid.
	InvalidMFA

	// PasswordExpired when password is expired.
	PasswordExpired

	// InvalidCredentials when credentials are invalid.
	InvalidCredentials

	// AlreadyExists when account already exists.
	AlreadyExists

	// Locked when account is locked.
	Locked

	// BannedPermanently when account is banned permanently.
	BannedPeramnently

	// Deleted when account is soft deleted.
	Deleted

	// PermissionDenied when account has no permission to perform the operation.
	PermissionDenied

	// Invalid when some invalid data is provided, or invalid operation is requested.
	// This is a generic error code, use more specific code if possible.
	Invalid = 9999
)
