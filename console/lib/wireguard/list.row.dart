import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;
import './edit.dart';
import './icon.checkmark.dart';

void _Noop(api.Wireguard up) {}

class ListRow extends StatelessWidget {
  final api.Wireguard current;
  final bool active;
  final Widget leading;
  final Widget trailing;

  final Future<void> Function()? onTap;
  final void Function(api.Wireguard upd) onChange;
  final void Function(api.Wireguard upd) onDelete;
  const ListRow(
    this.current, {
    super.key,
    this.active = false,
    this.leading = const SizedBox.shrink(),
    this.trailing = const SizedBox.shrink(),
    this.onChange = _Noop,
    this.onDelete = _Noop,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ds.TableRow(
      key: ValueKey(current.id),
      expanded: Edit(
        current,
        onChange: (upd) {
          return httpx
              .withRetry(
                () => api.wireguard.update(
                  upd,
                  options: [authn.AuthzCache.bearer(context)],
                ),
              )
              .then((resp) {
                if (resp.wireguard.default_5) {
                  return httpx.withRetry(
                    () => api.wireguard
                        .touch(
                          resp.wireguard.id,
                          options: [authn.AuthzCache.bearer(context)],
                        )
                        .then((_) => resp),
                  );
                }

                return Future.value(resp);
              })
              .then((resp) {
                onChange(resp.wireguard);
              });
        },
      ),
      [
        IconCheckmark(active, onTap: onTap),
        leading,
        Expanded(
          child: Text(
            current.description,
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        ),
        trailing,
      ],
    );
  }
}
