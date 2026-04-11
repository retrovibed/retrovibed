import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/httpx.dart' as httpx;

class CreateInlined extends StatefulWidget {
  final ds.AsyncVoidCallback onClose;
  final Future<meta.Profile> Function(meta.Profile)? onCreated;

  const CreateInlined({super.key, required this.onClose, this.onCreated});

  @override
  State<CreateInlined> createState() => _CreateInlinedState();
}

class _CreateInlinedState extends State<CreateInlined> {
  bool _loading = false;
  Widget _cause = ds.Error.zero;
  meta.Profile _profile = meta.Profile();
  meta.Token _token = meta.Token()..libraryRead = true;
  String _publicKey = '';

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _clearCause() => setState(() => _cause = ds.Error.zero);

  Future<void> _submit() {
    if (_publicKey.isEmpty) {
      setState(
        () =>
            _cause = ds.Error.text(
              "Public key is required",
              onTap: _clearCause,
            ),
      );
      return Future.value();
    }

    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final request =
        meta.ProfileCreateRequest()
          ..profile = _profile
          ..publicKey = _publicKey;

    return httpx
        .withRetry(
          () => meta.profiles
              .create(request, options: [authn.AuthzCache.bearer(context)])
              .then((v) {
                return httpx.withRetry(
                  () => meta.authz
                      .grant(
                        v.profile.id,
                        _token,
                        options: [authn.AuthzCache.bearer(context)],
                      )
                      .then((_) => v.profile),
                );
              }),
        )
        .then((profile) {
          setState(() => _loading = false);
          widget.onCreated?.call(profile);
          widget.onClose();
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: _clearCause);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      ds.Container(
        padding: defaults.padding,
        margin: defaults.margin,
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            profiles.Create(
              _profile,
              _publicKey,
              _token,
              onChange:
                  (profile, key, token) => setState(() {
                    _profile = profile;
                    _publicKey = key;
                    _token = token;
                  }),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              spacing: defaults.spacing,
              children: [
                ds.LoadingButton(Text("Cancel"), onPressed: widget.onClose),
                ds.LoadingButton(Text("Add"), onPressed: _submit),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
