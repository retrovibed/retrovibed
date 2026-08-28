import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;

class SettingsAudioSink extends StatefulWidget {
  const SettingsAudioSink({super.key, this.margin = EdgeInsets.zero, this.padding});
  final EdgeInsets margin;
  final EdgeInsets? padding;

  @override
  State<SettingsAudioSink> createState() => _SettingsAudioSinkState();
}

class _SettingsAudioSinkState extends State<SettingsAudioSink> with LoadingState {
  Widget unsupported(Object cause) => ds.Error.unavailable(
    cause,
    message: Text(
      "output only available on linux devices",
      style: TextStyle(color: Theme.of(context).disabledColor),
    ),
    decoration: ds.ErrorDecorations.info,
    color: Colors.transparent,
  );
  List<api.AudioSink> _sinks = [];
  String? _currentId;
  StreamSubscription<api.AudioSink>? _subscription;

  void _connect() {
    _subscription?.cancel();
    _sinks = [];

    api.sinks
        .listen(options: [authn.request(authn.AuthzCache.meta(context))])
        .then((stream) {
          _subscription = stream.listen(
            (v) {
              setState(() {
                _sinks = [..._sinks, v];
              });
            },
            onError: (cause) {
              setState(() {
                this.cause = ds.Error.unknown(cause, onTap: reseterr);
              });
            },
          );
        })
        .catchError((cause) {
          setState(() {
            this.cause = unsupported(cause);
          });
        }, test: httpx.ErrorsTest.unavailable)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });

    api.sinks
        .current(options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _currentId = v.sink.id;
            loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            this.cause = unsupported(cause);
          });
        }, test: httpx.ErrorsTest.unavailable)
        .catchError((cause) {
          setState(() {
            loading = false;
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  void _activate(String id) {
    setState(() => loading = true);

    api.sinks
        .activate(id, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _currentId = v.sink.id;
            loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            loading = false;
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_connect);
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      alignment: Alignment.topLeft,
      padding: widget.padding ?? defaults.padding,
      margin: widget.margin,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing / 2,
        children: [
          forms.Field(
            label: const Text("Output device"),
            cause: cause,
            input: DropdownButton<String>(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _currentId,
              items: _sinks.map((s) {
                return DropdownMenuItem<String>(
                  value: s.id,
                  child: Tooltip(
                    message: s.name,
                    child: Text(s.name, overflow: TextOverflow.ellipsis),
                  ),
                );
              }).toList(),
              onChanged: loading
                  ? null
                  : (String? v) {
                      if (v == null) return;
                      _activate(v);
                    },
            ),
          ),
        ],
      ),
    );
  }
}
