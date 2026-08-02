import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;

class ManualConfiguration extends StatefulWidget {
  final void Function()? retry;
  final void Function(api.Daemon) connect;
  final Alignment? alignment;

  ManualConfiguration({
    super.key,
    this.retry,
    required this.connect,
    this.alignment = Alignment.topCenter,
  });

  @override
  State<ManualConfiguration> createState() => _ManualConfigurationView();
}

class _ManualConfigurationView extends State<ManualConfiguration> with LoadingState {
  final String defaultLocalhost = httpx.localhost();
  String _hostname = '';

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      alignment: widget.alignment,
      padding: defaults.padding,
      ds.ErrorScreen(
        cause: cause,
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(
              label: SelectableText("hostname"),
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
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                if (widget.retry != null)
                  ds.LoadingButton(
                    Text("retry"),
                    onPressed: () async => widget.retry!(),
                  ),
                ds.LoadingButton(
                  Text("connect"),
                  onPressed: () {
                    return api.daemons
                        .create(
                          api.DaemonCreateRequest(
                            daemon: api.Daemon(
                              hostname: _hostname.isEmpty ? defaultLocalhost : _hostname,
                            ),
                          ),
                        )
                        .then((d) {
                          return widget.connect(d.daemon);
                        })
                        .catchError((error) {
                          setState(() {
                            cause = ds.Errors.httpauto(
                              error,
                              onTap: resetCause,
                            );
                          });
                        }, test: httpx.ErrorsTest.httpauto)
                        .catchError((error) {
                          setState(() {
                            cause = ds.Error.unknown(error, onTap: resetCause);
                          });
                        });
                  },
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
