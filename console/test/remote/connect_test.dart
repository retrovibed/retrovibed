import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/remote/connect.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  testWidgets('unmounting Connect does not throw during dispose', (tester) async {
    bool visible = true;
    late StateSetter setLocalState;

    await tester.pumpApp(
      StatefulBuilder(
        builder: (context, setState) {
          setLocalState = setState;
          return visible
              ? Connect(
                  search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
                  mode: ValueNotifier(media.SearchMode.library),
                  onModeChanged: (_) {},
                  downloading: const SizedBox(),
                  onDownloadingChanged: (_) {},
                )
              : const SizedBox();
        },
      ),
    );
    await tester.pump();

    setLocalState(() => visible = false);
    await tester.pump();

    expect(tester.takeException(), isNull);
  });
}
