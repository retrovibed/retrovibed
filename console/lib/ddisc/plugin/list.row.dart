import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import './environment.editor.dart';

class ListRow extends StatelessWidget {
  final api.Plugin current;
  final void Function(api.Plugin deleted) onDelete;

  const ListRow(this.current, {super.key, this.onDelete = ds.fnNoop});

  @override
  Widget build(BuildContext context) {
    return ds.TableRow(
      key: ValueKey(current.id),
      expanded: EnvironmentEditor.future(
        current.id,
        api.environment.get(
          current.id,
          options: [authn.request(authn.AuthzCache.meta(context))],
        ),
        onChange: (content) {
          httpx
              .withRetry(
                () => api.environment.update(
                  current.id,
                  content,
                  options: [authn.request(authn.AuthzCache.meta(context))],
                ),
              )
              .catchError((cause) {
                print("failed to update plugin environment ${cause}");
                return content;
              });
        },
      ),
      [
        Expanded(
          child: Text(
            current.name,
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        ),
        Text(ds.bytesx(current.size.toInt()).toIEC600272Format()),
        IconButton(
          icon: const Icon(Icons.delete_outline),
          onPressed: () {
            final modal = ds.modals.of(context);
            modal?.push(
              ds.Confirmation.yesNo(
                content: Text('Delete ${current.name}?'),
                onCancel: (_) => modal.push(null),
                onConfirm: (_) {
                  api.plugins
                      .delete(
                        current.id,
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      )
                      .then((_) {
                        modal.push(null);
                        onDelete(current);
                      });
                },
              ),
            );
          },
        ),
      ],
    );
  }
}
