import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;

class ManualConfiguration extends StatefulWidget {
  final void Function()? retry;
  final void Function(api.Daemon) onConnected;
  final Future<api.Daemon> Function(api.Daemon) apiconnect;
  final Future<api.DaemonCreateResponse> Function(api.DaemonCreateRequest req) apicreate;
  final Alignment? alignment;

  ManualConfiguration(
    this.onConnected, {
    super.key,
    this.retry,
    this.apiconnect = api.daemons.connectable,
    this.apicreate = api.daemons.create,
    this.alignment = Alignment.topCenter,
  });

  @override
  State<ManualConfiguration> createState() => _ManualConfigurationView();
}

class _ManualConfigurationView extends State<ManualConfiguration> with LoadingState {
  final String defaultLocalhost = httpx.localhost();
  String _hostname = '';

  @override
  void initState() {
    super.initState();
    loading = false;
  }

  Future<void> _connect() async {
    setState(() {
      loading = true;
    });

    final daemon = api.Daemon(
      hostname: _hostname.isEmpty ? defaultLocalhost : _hostname,
    );

    return httpx
        .withRetry(
          () async {
            return widget.apiconnect(daemon).then((d) => _create(d));
          },
          maxRetries: 2,
          backoff: (attempt) => Duration(milliseconds: 200),
          checks: [
            httpx.RetryChecks.ratelimited,
            httpx.RetryChecks.badgateway,
            httpx.RetryChecks.unavailable,
          ],
        )
        .catchError((error) async {
          setState(() {
            loading = false;
          });

          final proceed = await ds.modals.asyncfn<bool>(context, (completion) {
            return ds.Confirmation(
              constraints: BoxConstraints(maxWidth: 512),
              content: ds.Error.offline(
                error,
                constraints: BoxConstraints(minHeight: 128),
              ),
              confirmation: Text("Create Anyway"),
              cancellation: Text("Cancel"),
              onConfirm: (_) => completion.complete(true),
              onCancel: (_) => completion.complete(false),
            );
          });

          if (proceed) {
            return _create(daemon);
          }

          return;
        }, test: ds.ErrorTests.offline)
        .catchError((error) async {
          setState(() {
            loading = false;
          });

          final proceed = await ds.modals.asyncfn<bool>(context, (completion) {
            return ds.Confirmation(
              constraints: BoxConstraints(maxWidth: 512),
              content: ds.Error.dnsresolution(
                error,
                constraints: BoxConstraints(minHeight: 128),
              ),
              confirmation: Text("Create Anyway"),
              cancellation: Text("Cancel"),
              onConfirm: (_) => completion.complete(true),
              onCancel: (_) => completion.complete(false),
            );
          });

          if (proceed) {
            return _create(daemon);
          }

          return;
        }, test: ds.ErrorTests.dnsresolution)
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(error, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(error, onTap: reseterr);
          });
        });
  }

  Future<void> _create(api.Daemon daemon) {
    return widget
        .apicreate(api.DaemonCreateRequest(daemon: daemon))
        .then((response) {
          setState(() {
            loading = false;
          });
          widget.onConnected(response.daemon);
        })
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(error, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      alignment: widget.alignment,
      padding: defaults.padding,
      decoration: BoxDecoration(
        border: defaults.border,
        borderRadius: defaults.borderRadius,
      ),
      cause: cause,
      loading: loading,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            label: Text("hostname"),
            input: TextFormField(
              autofocus: true,
              decoration: InputDecoration(
                hintText: defaultLocalhost,
                helperText: "hostname and port for the retrovibed instance",
              ),
              onChanged: (v) {
                setState(() {
                  _hostname = v;
                });
              },
              onFieldSubmitted: (_) => _connect(),
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            spacing: defaults.spacing,
            children: [
              if (widget.retry != null)
                ds.LoadingButton(
                  Text("retry"),
                  border: defaults.border,
                  onPressed: () async => widget.retry!(),
                ),
              ds.LoadingButton(
                Text("connect"),
                border: defaults.border,
                onPressed: _connect,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
