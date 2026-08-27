import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show TextInput;
import 'package:flutter/foundation.dart' as foundation;
import 'package:window_manager/window_manager.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/retrovibed.dart' as retro;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'developer.mode.dart';

class Login extends StatefulWidget {
  final Widget child;
  final String Function() publicKey;
  final String Function(String) seed;
  final bool Function() guest;
  final Future<void> Function() authenticated;

  const Login(
    this.child, {
    super.key,
    this.publicKey = retro.public_key,
    this.seed = retro.seed,
    this.guest = retro.guest,
    this.authenticated = _noop,
  });

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

class _LoginState extends State<Login> {
  Widget _cause = ds.Error.zero;
  bool _isObscured = true;
  bool _hasKey = false;
  bool _acceptedTos = false;
  bool _register = false;
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
    _register = !widget.publicKey().isNotEmpty;
    _checkKey();
  }

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _logout() {
    retro.unseed();
    setState(() {
      _hasKey = false;
    });
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _checkKey() {
    if (!widget.publicKey().isNotEmpty) return;
    if (_hasKey) return;

    widget
        .authenticated()
        .then((_) {
          setState(() {
            _hasKey = true;
          });
        })
        .catchError((e) {
          setState(() {
            _hasKey = false;
            _cause = ds.Error.unknown(e, onTap: _reseterr);
          });
        });
  }

  Future<void> _seed() async {
    _reseterr();
    if (_register && _password != _confirm) {
      setState(() {
        _cause = ds.Error.text("passwords do not match", onTap: _reseterr);
      });
      return;
    }

    return Future.microtask(() {
          final err = widget.seed("${_username}:${_password}");
          if (err.isNotEmpty) {
            return Future.error(err);
          }
          return Future.value();
        })
        .then((_ignored) {
          TextInput.finishAutofillContext();
          return _ignored;
        })
        .catchError((cause) {
          print(cause);
          setState(() {
            _cause = ds.Error.text("login failed", onTap: _reseterr);
          });
        })
        .then((_) => _checkKey());
  }

  Future<void> _guestLogin() async {
    _reseterr();
    if (!widget.guest()) {
      setState(() {
        _cause = ds.Error.text("guest login failed", onTap: _reseterr);
      });
      return;
    }
    _checkKey();
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
            ds.Loading(
              cause: _cause,
              ds.Container(
                padding: defaults.padding,
                margin: defaults.margin,
                constraints: BoxConstraints(maxWidth: 375),

                SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    spacing: defaults.spacing,
                    children: [
                      SizedBox(
                        width: double.infinity,
                        child: Stack(
                          alignment: Alignment.center,
                          children: [
                            Text(
                              'Welcome to Retrovibed',
                              style: Theme.of(context).textTheme.headlineSmall,
                              textAlign: TextAlign.center,
                            ),
                            Positioned(
                              left: 0,
                              child: ds.LoadingIconButton.guest(
                                tooltip: "continue as guest",
                                onPressed: _guestLogin,
                              ),
                            ),
                            Positioned(
                              right: 0,
                              child: ds.LoadingIconButton.close(
                                tooltip: "exit application",
                                onPressed: windowManager.close,
                              ),
                            ),
                          ],
                        ),
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
                            TextFormField(
                              initialValue: _username,
                              autofillHints: const [AutofillHints.username, AutofillHints.email],
                              keyboardType: TextInputType.emailAddress,
                              decoration: InputDecoration(hintText: 'email'),
                              onChanged: (v) => setState(() => _username = v),
                              onFieldSubmitted: (_) => _seed(),
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
                              onFieldSubmitted: (_) => _seed(),
                            ),
                            Visibility(
                              visible: _register,
                              child: TextFormField(
                                obscureText: _isObscured,
                                autofillHints: const [AutofillHints.password],
                                decoration: InputDecoration(
                                  hintText: 'confirm password',
                                  suffixIcon: obscureicon,
                                ),
                                onChanged: (v) => setState(() => _confirm = v),
                                onFieldSubmitted: (_) => _seed(),
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
                        disabled: _username.isEmpty || _password.isEmpty || !_acceptedTos,
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
