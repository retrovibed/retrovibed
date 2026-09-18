import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;

import './api.dart' as api;

class DownloadDisplay extends StatelessWidget {
  final api.Download current;
  final List<Widget> trailing;
  final Future<void> Function()? onReset;
  final Future<void> Function()? onDelete;
  final Future<void> Function(api.Download)? onVerify;
  const DownloadDisplay(
    this.current, {
    super.key,
    this.onReset,
    this.onDelete,
    this.onVerify,
    this.trailing = const [],
  });

  static Widget fromID(
    String id, {
    Key? key,
    List<Widget> trailing = const [],
    Future<void> Function()? onReset,
    Future<void> Function()? onDelete,
    Future<void> Function(api.Download)? onVerify,
    Future<api.DownloadMetadataResponse> Function(String id, {List<httpx.Option> options}) get = api.discovered.get,
  }) {
    return Builder(
      key: key,
      builder: (context) {
        return FutureBuilder<api.Download>(
          initialData: api.Download.create(),
          // an id with no backing record (e.g. deleted between load and
          // fetch) resolves to an empty Download rather than an error, and
          // the display renders disabled.
          future: get(
            id,
            options: [authn.request(authn.AuthzCache.meta(context))],
          ).then((v) => v.download).catchError((_) => api.Download(), test: httpx.ErrorsTest.err404),
          builder: (BuildContext ctx, AsyncSnapshot<api.Download> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              DownloadDisplay(
                snapshot.data ?? api.Download(),
                trailing: trailing,
                onReset: onReset,
                onDelete: onDelete,
                onVerify: onVerify,
              ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Opacity(
      opacity: current.hasMedia() ? 1.0 : 0.4,
      child: IgnorePointer(
        ignoring: !current.hasMedia(),
        child: ds.Container(
          padding: defaults.padding,
          margin: defaults.margin.copyWith(bottom: 0),
          Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              forms.Field(
                label: Text("id"),
                input: Text(
                  current.media.id,
                  overflow: TextOverflow.ellipsis,
                  maxLines: 1,
                ),
                trailing: [
                  if (onVerify != null)
                    ds.LoadingIconButton(
                      onPressed: () => onVerify!(current),
                      tooltip: "verify data",
                      icon: Icon(Icons.fact_check),
                    ),
                  if (onReset != null)
                    ds.LoadingIconButton.refresh(onPressed: onReset!, tooltip: "clear data from disk keeps metadata"),
                  if (onDelete != null)
                    ds.LoadingIconButton.delete(onPressed: onDelete!, tooltip: "permanently delete this media"),
                ],
              ),
              forms.Field(
                label: Text("description"),
                input: Text(current.media.description),
              ),
              forms.Field(label: Text("path"), input: Text(current.path)),
              forms.Field(label: Text("bytes"), input: ds.Bytes(current.bytes)),
              forms.Field(
                label: Text("paused"),
                input: ds.Timestamp.iso8601(current.pausedAt),
              ),
              Transform.translate(
                child: forms.Checkbox(
                  Text("distributing"),
                  value: current.distributing,
                ),
                offset: Offset(-10.0, 0.0),
              ),
              ...trailing,
            ],
          ),
        ),
      ),
    );
  }
}
