import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'metadata.typography.dart' as typography;
import 'metadata.icons.dart' as icons;

class MediaEdit extends StatelessWidget {
  final media.Media current;
  final Function(Future<media.Media>) onChange;
  final Widget closable;
  final Widget deletable;
  final EdgeInsets? padding;
  MediaEdit({
    super.key,
    required this.current,
    required this.onChange,
    this.closable = ds.Empty,
    this.deletable = ds.Empty,
    this.padding,
  });

  Widget archive(BuildContext context, String uid) {
    final display = typography.archived(uid);
    final icon = icons.archived_trash(uid);

    if (uuidx.isMax(uuidx.fromString(uid))) {
      return Row(
        children: [
          display,
          Spacer(),
          ds.LoadingIconButton(
            onPressed: () {
              return media.media
                  .update(
                    current.id,
                    current..archiveId = uuidx.min(),
                    options: [authn.request(authn.AuthzCache.meta(context))],
                  )
                  .then((v) => onChange(Future.value(v.media)));
            },
            icon: icon,
          ),
        ],
      );
    }

    if (uuidx.isMin(uuidx.fromString(uid))) {
      return Row(
        children: [
          display,
          Spacer(),
          ds.LoadingIconButton(
            onPressed: () {
              return media.media
                  .update(
                    current.id,
                    current..archiveId = uuidx.max(),
                    options: [authn.request(authn.AuthzCache.meta(context))],
                  )
                  .then((v) => onChange(Future.value(v.media)));
            },
            icon: icon,
          ),
        ],
      );
    }

    return Row(
      children: [
        display,
        Spacer(),
        ds.LoadingIconButton(
          onPressed: () => media.media
              .unarchive(
                current.archiveId,
                options: [authn.Authenticated.bearer(context)],
              )
              .then(
                (v) => media.media.update(
                  current.id,
                  current..archiveId = uuidx.min(),
                  options: [authn.request(authn.AuthzCache.meta(context))],
                ),
              )
              .then((v) => onChange(Future.value(v.media))),
          icon: icon,
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      padding: padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(
            label: Text("id"),
            input: Text(current.id),
            trailing: [deletable, closable],
          ),
          forms.Field(
            label: Text("description"),
            input: TextFormField(
              initialValue: current.description,
              onChanged: (v) => onChange(Future.value(current..description = v)),
            ),
          ),
          forms.Field(
            label: Text("created"),
            input: ds.Timestamp.iso8601(current.createdAt),
          ),
          forms.Field(label: Text("mimetype"), input: Text(current.mimetype)),
          forms.Field(
            label: Text("sharing"),
            input: typography.sharing(current.torrentId),
          ),
          Visibility(
            visible: (authn.AuthzCache.of(context).meta.current.token.archiveUpload.toInt()) > 0,
            child: forms.Field(
              label: Text("archived"),
              input: archive(context, current.archiveId),
            ),
          ),
        ],
      ),
    );
  }
}
