import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/typography.dart' as typography;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;

import './api.dart' as api;

class DownloadDisplay extends StatelessWidget {
  final api.Download current;
  final Color? background;
  final List<Widget> trailing;
  final Future<void> Function()? onTap;
  final Future<void> Function(api.Download)? onVerify;
  const DownloadDisplay(
    this.current, {
    super.key,
    this.onTap,
    this.onVerify,
    this.background,
    this.trailing = const [],
  });

  static Widget fromID(
    String id, {
    Key? key,
    Color? background,
    List<Widget> trailing = const [],
    Future<void> Function()? onTap,
    Future<void> Function(api.Download)? onVerify,
    Future<api.DownloadMetadataResponse> Function(String id, {List<httpx.Option> options}) get = api.discovered.get,
  }) {
    return Builder(
      key: key,
      builder: (context) {
        return FutureBuilder<api.Download>(
          initialData: api.Download.create(),
          future: get(id, options: [authn.request(authn.AuthzCache.meta(context))]).then((v) => v.download),
          builder: (BuildContext ctx, AsyncSnapshot<api.Download> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              snapshot.data == null
                  ? SizedBox()
                  : DownloadDisplay(
                    snapshot.data!,
                    background: background,
                    trailing: trailing,
                    onTap: onTap,
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
    return ds.Container(
      decoration: BoxDecoration(color: background),
      padding: defaults.padding,
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
              if (onTap != null)
                ds.LoadingIconButton.delete(onPressed: onTap!, tooltip: "clear data from disk keeps metadata"),
            ],
          ),
          forms.Field(
            label: Text("description"),
            input: Text(current.media.description),
          ),
          forms.Field(label: Text("path"), input: Text(current.path)),
          forms.Field(
            label: Text("bytes"),
            input: typography.Bytes(current.bytes),
          ),
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
    );
  }
}
