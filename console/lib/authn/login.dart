import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/retrovibed.dart' as retro;

class Login extends StatefulWidget {
  final Widget child;
  final String Function() publicKey;
  final String Function(String) seed;
  final Future<void> Function() authenticated;

  const Login(
    this.child, {
    super.key,
    this.publicKey = retro.public_key,
    this.seed = retro.seed,
    this.authenticated = _noop,
  });

  static Future<void> _noop() => Future.value();

  static void logout(BuildContext context) {
    context.findAncestorStateOfType<_LoginState>()?._logout();
  }

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  Widget _cause = ds.Error.zero;
  bool _isObscured = true;
  bool _hasKey = false;
  bool _acceptedTos = false;
  String _username = '';
  String _password = '';

  @override
  void initState() {
    super.initState();
    _checkKey();
  }

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _logout() {
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
    final err = widget.seed("${_username}:${_password}");
    if (err.isNotEmpty) {
      setState(() {
        _cause = ds.Error.text("login failed", onTap: _reseterr);
      });
      return;
    }
    _checkKey();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    if (_hasKey) return widget.child;

    return ds.Masked(
      alignment: Alignment.center,
      Center(
        child: ds.Container(
          padding: defaults.padding,
          margin: defaults.margin,
          constraints: BoxConstraints(maxWidth: 375),
          ds.Loading(
            cause: _cause,
            Column(
              mainAxisSize: MainAxisSize.min,
              spacing: defaults.spacing,
              children: [
                Text(
                  'Welcome to Retrovibed',
                  style: Theme.of(context).textTheme.headlineSmall,
                  textAlign: TextAlign.center,
                ),
                Text(
                  'setup your device',
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
                const SizedBox(height: 16),
                TextFormField(
                  decoration: InputDecoration(hintText: 'email'),
                  onChanged: (v) => setState(() => _username = v),
                  onFieldSubmitted: (_) => _seed(),
                ),
                TextFormField(
                  obscureText: _isObscured,
                  decoration: InputDecoration(
                    hintText: 'password',
                    suffixIcon: IconButton(
                      icon: Icon(
                        _isObscured ? Icons.visibility : Icons.visibility_off,
                      ),
                      onPressed: () {
                        setState(() {
                          _isObscured = !_isObscured;
                        });
                      },
                    ),
                  ),
                  onChanged: (v) => setState(() => _password = v),
                  onFieldSubmitted: (_) => _seed(),
                ),
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Checkbox(
                      value: _acceptedTos,
                      onChanged: (v) => setState(() => _acceptedTos = v ?? false),
                    ),
                    Flexible(
                      child: Text.rich(
                        textAlign: TextAlign.center,
                        TextSpan(
                          text: 'By continuing you accept the ',
                          children: [
                            ds.Hyperlink.inline(
                              'terms of service',
                              url: 'https://retrovibe.space/terms',
                            ),
                          ],
                        ),
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                    ),
                  ],
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
    );
  }
}
