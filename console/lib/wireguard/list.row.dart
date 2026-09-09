import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;
import 'edit.dart';

class ListRow extends StatelessWidget {
  final api.Wireguard current;
  final List<Widget> leading;
  final List<Widget> trailing;

  final Future<void> Function()? onTap;
  final Future<void> Function(api.Wireguard o, api.Wireguard upd) onChange;
  final Future<void> Function(api.Wireguard o) onDelete;
  const ListRow(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
    this.onChange = ds.fnAsyncNoopOnChangeV2,
    this.onDelete = ds.fnAsyncNoopOnDelete,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ds.TableRow(
      key: ValueKey(current.id),
      expanded: Edit(
        current,
        onChange: (o, upd) {
          return httpx
              .withRetry(
                () => api.wireguard.update(
                  upd,
                  options: [authn.request(authn.AuthzCache.meta(context))],
                ),
              )
              .then((resp) {
                if (resp.wireguard.nettype == api.WireguardNettype.UNSPECIFIED) {
                  return Future.value(resp);
                }

                return httpx.withRetry(
                  () => api.wireguard
                      .touch(
                        resp.wireguard.id,
                        resp.wireguard.nettype,
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      )
                      .then((_) => resp),
                );
              })
              .then((resp) {
                onChange(current, resp.wireguard);
              });
        },
        onDelete: onDelete,
      ),
      [
        ...leading,
        Expanded(
          child: Text(
            current.description,
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        ),
        ...trailing,
      ],
    );
  }
}
