import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;

// deleting a directory deletes what it holds, which is not recoverable from this screen,
// so the user is told before it happens rather than after.
Future<void> confirmremove(
  BuildContext context,
  media.Media current, {
  required api.FnFilesystemDelete apiremove,
  required void Function(media.Media? upd) onChange,
}) {
  return ds.modals.asyncfn<void>(
    context,
    ds.Confirmation.dangerous(
      content: Text(
        "Delete ${current.description}? Everything inside it is removed from your library too.",
      ),
      onConfirm: (ctx) => httpx
          .withRetry(
            () => apiremove(current.id, options: [authn.request(authn.AuthzCache.meta(ctx))]),
          )
          .then((_) => onChange(null)),
    ),
  );
}

// a directory is only a container, so it carries its type and the general actions; files
// add the metadata the listing row has no room for.
class FilesystemDetails extends StatelessWidget {
  final media.Media current;
  // null signals the entry was removed.
  final void Function(media.Media? upd) onChange;
  final api.FnFilesystemDelete apiremove;

  const FilesystemDetails(
    this.current, {
    super.key,
    required this.onChange,
    this.apiremove = api.filesystem.delete,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final directory = current.mimetype == mimex.directory;

    return forms.Container(
      padding: defaults.padding,
      margin: defaults.margin.copyWith(bottom: 0),
      Column(
        mainAxisSize: MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          forms.Field(
            label: const Text("type"),
            input: Text(current.mimetype),
            trailing: [
              // a directory has no content of its own to fetch.
              if (!directory)
                ds.LoadingIconButton(
                  onPressed: media.DownloadAction(context, current),
                  icon: const Icon(Icons.download),
                  tooltip: "download",
                ),
              ds.LoadingIconButton.delete(
                color: defaults.danger,
                onPressed: () => confirmremove(context, current, apiremove: apiremove, onChange: onChange),
                tooltip: "permanently delete",
              ),
            ],
          ),
          if (!directory) ...[
            forms.Field(
              label: const Text("id"),
              input: Text(current.id, overflow: TextOverflow.ellipsis, maxLines: 1),
            ),
            forms.Field(label: const Text("created"), input: ds.Timestamp.iso8601(current.createdAt)),
            forms.Field(label: const Text("updated"), input: ds.Timestamp.iso8601(current.updatedAt)),
          ],
        ],
      ),
    );
  }
}
