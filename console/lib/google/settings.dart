import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/google/api.dart' as api;
import 'package:url_launcher/url_launcher.dart';

class Settings extends StatefulWidget {
  const Settings({super.key});

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> with ds.LoadingState {
  Timer? _poll;
  api.YouTubeStatus _youtube = api.YouTubeStatus();

  @override
  void initState() {
    super.initState();
    ds.postframe(() => _fetch());
  }

  @override
  void dispose() {
    _poll?.cancel();
    super.dispose();
  }

  void _fetch() {
    setState(() {
      loading = true;
      cause = ds.Error.zero;
    });

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    httpx
        .withRetry(() => api.YouTube.status(options: auth))
        .then((v) {
          if (v.linked) _poll?.cancel();
          setState(() {
            loading = false;
            _youtube = v;
          });
        })
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(e, onTap: _fetch);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    return ds.Loading(
      loading: loading,
      cause: cause,
      Padding(
        padding: defaults.padding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          spacing: defaults.spacing,
          children: [
            Text("YouTube", style: theme.textTheme.titleMedium),
            Text(
              _youtube.linked ? "Connected" : "Not connected",
              style: theme.textTheme.bodySmall,
            ),
            SizedBox(
              width: double.infinity,
              child:
                  _youtube.linked
                      ? ds.LoadingButton(
                        Text("Unlink"),
                        onPressed: () {
                          return httpx
                              .withRetry(
                                () => api.YouTube.unlink(options: auth),
                              )
                              .then((_) => _fetch())
                              .catchError((e) {
                                setState(() {
                                  cause = ds.Error.unknown(e, onTap: _fetch);
                                });
                              });
                        },
                      )
                      : OutlinedButton(
                        onPressed: () {
                          authn.otp(options: [authn.DeeppoolAuthzCache.bearer(context)]).then((session) {
                            launchUrl(
                              api.YouTube.authUri(token: session.token),
                            );
                            _poll?.cancel();
                            _poll = Timer.periodic(
                              const Duration(seconds: 3),
                              (_) => _fetch(),
                            );
                          });
                        },
                        child: Text("Link YouTube"),
                      ),
            ),
          ],
        ),
      ),
    );
  }
}
