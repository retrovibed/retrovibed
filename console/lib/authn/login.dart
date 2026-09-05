import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show TextInput;
import 'package:flutter/foundation.dart' as foundation;
import 'package:retrovibed/windowx.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/retrovibed.dart' as retro;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'developer.mode.dart';

class Login extends StatefulWidget {
  final Widget child;
  final String Function() publicKey;
  final Future<void> Function(String, String) seed;
  final Future<void> Function() authenticated;
  final Future<void> Function() unseed;
  final Future<void> Function() guest;
  final WindowManagerX windowManager;

  Login(
    this.child, {
    super.key,
    this.publicKey = retro.public_key,
    this.seed = retro.seed,
    this.unseed = retro.unseed,
    this.guest = retro.guest,
    this.authenticated = _noop,
    WindowManagerX? windowManager,
  }) : windowManager = windowManager ?? windowx;

  static Future<void> _noop() => Future.value();

  static void logout(BuildContext context) {
    context.findAncestorStateOfType<_LoginState>()?._logout();
  }

  static _LoginCachedData cached(BuildContext context) {
    return context.dependOnInheritedWidgetOfExactType<_LoginCachedData>() ?? _LoginCachedData.empty;
  }

  static _LoginState? of(BuildContext context) {
    return context.findAncestorStateOfType<_LoginState>();
  }

  @override
  State<Login> createState() => _LoginState();
}

class _LoginCachedData extends InheritedWidget {
  final DeveloperMode flags;

  const _LoginCachedData({required this.flags, required super.child});

  static final empty = _LoginCachedData(
    flags: DeveloperMode(),
    child: const SizedBox(),
  );

  @override
  bool updateShouldNotify(_LoginCachedData old) => flags != old.flags;
}

class _LoginState extends State<Login> with ds.LoadingState {
  bool _isObscured = true;
  bool _register = false;
  bool _hasKey = false;
  bool _acceptedTos = false;
  String _username = '';
  String _password = '';
  String _confirm = '';
  DeveloperMode flags = DeveloperMode(
    alpha: foundation.kDebugMode,
    recommendations: true,
    releases: true,
    debug: foundation.kDebugMode,
    subscription: !(Platform.isAndroid || Platform.isIOS),
  );

  @override
  void initState() {
    super.initState();
    _checkKey().then((_) {
      setState(() {
        loading = false;
      });
    });
  }

  Future<void> _logout() {
    return widget
        .unseed()
        .then((_) {
          setState(() {
            loading = false;
            _hasKey = false;
            _username = '';
            _password = '';
          });
        })
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  Future<void> _checkKey() {
    if (widget.publicKey().isEmpty) return Future.sync(() {});
    if (_hasKey) return Future.sync(() {});

    return widget
        .authenticated()
        .then((_) {
          setState(() {
            loading = false;
            _hasKey = widget.publicKey().isNotEmpty;
          });
        })
        .catchError((e) {
          _logout().then((_) {
            setState(() {
              loading = false;
              cause = ds.Error.unknown(e, onTap: reseterr);
            });
          });
        });
  }

  Future<void> _seed() async {
    reseterr();

    // prevent resubmissions while running
    if (loading) return;
    setState(() {
      loading = true;
    });

    return widget
        .seed(_username, _password)
        .then((_ignored) {
          TextInput.finishAutofillContext();
          return _ignored;
        })
        .catchError((e) {
          print(e);
          setState(() {
            cause = ds.Error.text("login failed", onTap: reseterr);
          });
        })
        .then((_) => _checkKey())
        .whenComplete(() {
          setState(() {
            loading = false;
          });
        });
  }

  Future<void> _guestLogin() async {
    reseterr();
    return widget
        .guest()
        .then((_) {
          _checkKey();
        })
        .catchError((e) {
          print(e);
          setState(() {
            cause = ds.Error.text("guest login failed", onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    if (_hasKey) return _LoginCachedData(flags: flags, child: widget.child);
    final obscureicon = IconButton(
      icon: Icon(
        _isObscured ? Icons.visibility : Icons.visibility_off,
      ),
      onPressed: () {
        setState(() {
          _isObscured = !_isObscured;
        });
      },
    );

    return _LoginCachedData(
      flags: flags,
      child: ds.Masked(
        alignment: Alignment.center,
        modals.Node(
          ds.HelpGlobal(
            ds.LoadingBoundary(
              loading: loading,
              cause: cause,
              ds.Container(
                padding: defaults.padding,
                margin: defaults.margin,
                constraints: BoxConstraints(maxWidth: 375),
                SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    spacing: defaults.spacing,
                    children: [
                      Row(
                        children: [
                          ds.LoadingIconButton.guest(
                            tooltip: "continue as guest",
                            onPressed: _guestLogin,
                          ),
                          Expanded(
                            child: Text(
                              'Welcome to Retrovibed',
                              style: Theme.of(context).textTheme.headlineSmall,
                              textAlign: TextAlign.center,
                            ),
                          ),
                          ds.LoadingIconButton.close(
                            tooltip: "exit application",
                            onPressed: widget.windowManager.close,
                          ),
                        ],
                      ),
                      Text(
                        'setup your device',
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                      AutofillGroup(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          spacing: defaults.spacing,
                          children: [
                            SizedBox.square(dimension: 32),
                            TextFormField(
                              initialValue: _username,
                              autofillHints: const [AutofillHints.username, AutofillHints.email],
                              keyboardType: TextInputType.emailAddress,
                              decoration: InputDecoration(hintText: 'email'),
                              onChanged: (v) => setState(() => _username = v),
                              validator: (v) => (v ?? "").isEmpty ? "email cannot be empty" : null,
                            ),
                            TextFormField(
                              initialValue: _password,
                              obscureText: _isObscured,
                              autofillHints: const [AutofillHints.password],
                              decoration: InputDecoration(
                                hintText: 'password',
                                suffixIcon: obscureicon,
                              ),
                              onChanged: (v) => setState(() => _password = v),
                              validator: (v) => (v ?? "").isEmpty ? "password cannot be empty" : null,
                            ),
                            Visibility(
                              visible: _register,
                              maintainAnimation: true,
                              maintainSize: true,
                              maintainState: true,
                              child: TextFormField(
                                obscureText: _isObscured,
                                autofillHints: const [AutofillHints.password],
                                decoration: InputDecoration(
                                  hintText: 'confirm password',
                                  suffixIcon: obscureicon,
                                ),
                                onChanged: (v) => setState(() => _confirm = v),
                                validator: (v) => (v ?? "").isEmpty && v == _password ? "passwords must match" : null,
                              ),
                            ),
                          ],
                        ),
                      ),
                      forms.Checkbox(
                        Text('register a new account'),
                        value: _register,
                        onChanged: (v) => setState(() => _register = v ?? false),
                      ),
                      forms.Checkbox(
                        Text.rich(
                          TextSpan(
                            text: 'By continuing you accept the ',
                            children: [
                              ds.Hyperlink.inline(
                                'terms of service',
                                url: 'https://retrovibe.space/terms',
                              ),
                            ],
                          ),
                        ),
                        value: _acceptedTos,
                        onChanged: (v) => setState(() => _acceptedTos = v ?? false),
                      ),
                      ds.LoadingButton(
                        const Text('Login'),
                        onPressed: _seed,
                        disabled:
                            loading ||
                            _username.isEmpty ||
                            _password.isEmpty ||
                            (_register && (_password != _confirm)) ||
                            !_acceptedTos,
                      ),
                    ],
                  ),
                ),
              ),
            ),
            ds.Hint.multiline([
              Text(
                'Retrovibed takes personal privacy seriously. The information you enter here is used to seed cryptographic primatives and never leaves your device.',
              ),
            ]),
          ),
        ),
      ),
    );
  }
}
