// struct ListNode {
//     int val;
//     ListNode *next;
//     ListNode(int x) : val(x), next(NULL) {}
// };

#include <map>
using namespace std;
class Solution {
public:
    ListNode *getIntersectionNode(ListNode *headA, ListNode *headB) {
        map<ListNode*, bool> nodeMap;

        //List Aを先頭から見ていき、存在するノードをtrueと一緒にmapに入れる
        while(headA){
            nodeMap[headA] = true;
            headA = headA->next;
        }

        while (headB) {
            if (nodeMap.find(headB) != nodeMap.end()){//List Bの今見ているノードがnodeMapに見つかったら
                return headB;
            }
            headB = headB->next;
        }

        return nullptr;
    }
};