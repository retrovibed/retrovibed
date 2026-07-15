import 'package:retrovibed/media/media.pb.dart';

class MediaSearchState {
  MediaSearchRequest next;
  int count;

  MediaSearchState({required this.next, this.count = 0});
}
