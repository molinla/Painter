<script setup lang="ts">
//base
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useMessage, NAvatar, NIcon, NModal, NCard, NInput, NButton, NInputNumber, NUpload, type UploadFileInfo, type UploadCustomRequestOptions } from 'naive-ui'
//store
import { useUserStore } from '@/stores/user';
import { UserSelfFull, UserUpdate, UserHeaderFieldUpdate } from '@/apis/api_user';
//utils
import type { SelfItem } from '@/utils/interface';
import { dateDiff, dateToString } from "@/utils/timeToStr"

//icons
import { Add } from '@vicons/carbon'
import { KeyboardArrowRightFilled } from '@vicons/material';
import { LogoutOutlined } from '@vicons/antd';
import axios from 'axios';
//ref
const Router = useRouter()
const UserStore = useUserStore()
const Message = useMessage()
const { loginStatus } = storeToRefs(UserStore)
const showModal = ref(false)
const showAvatarModal = ref(false)
const userInfo = ref<SelfItem>({
    ArticleNumber: 0,
    ArticleList: [],
    CollectionNumber: 0,
    FollowingNumber: 0,
    FollowerNumber: 0,
    UserInfo: {
        ID: 0,
        Email: '',
        NickName: '',
        HeaderField: '',
        CreatedAt: '',
        LastLogin: '',
        UserName: '',
        PhoneNum: 0,
        UserGroup: 0
    },
    Following: false,
    TotalCount: 0
})
const fileList = ref<UploadFileInfo[]>([
    {
        id: '1',
        name: 'avatar',
        status: 'finished',
        url: userInfo.value.UserInfo.HeaderField
    },
])
//fn
onMounted(async () => {
    if (!loginStatus.value) {
        Router.push({ path: '/login' })
        Message.warning("账号还未登录，请登录再试！")
    }
    let res = await UserStore.loadSelfData()
    if (res.ok) {
        let info = await UserSelfFull()
        userInfo.value = info.data
        fileList.value[0].url = info.data.UserInfo.HeaderField
    } else {
        Router.push({ path: '/login' })
        Message.warning("账号还未登录，请登录再试！")
    }
})

const goArticleList = (id: number) => {
    Router.push({ path: "/articlelist", query: { type: 3, id: id } })
}

const goFollowing = (id: number) => {
    Router.push({ path: "/follow", query: { type: 1, id: id } })
}

const goFollower = (id: number) => {
    Router.push({ path: "/follow", query: { type: 2, id: id } })
}

const goCollections = (id: number) => {
    Router.push({ path: "/articlelist", query: { type: 4, id: id } })
}

const goEditArticle = (type: number, id: number) => {
    Router.push({ path: "/editarticle", query: { type: type, id: id } })
}

const doLogout = () => {
    UserStore.logout()
    Router.push({ path: "/login" })
}

const update = async () => {
    await UserUpdate({
        Name: userInfo.value.UserInfo.UserName,
        Email: userInfo.value.UserInfo.Email,
        PhoneNum: userInfo.value.UserInfo.PhoneNum as number,
        NickName: userInfo.value.UserInfo.NickName,
    })
    showModal.value = !showModal.value
}

const showChangeAvatarModal = async () => {
    showAvatarModal.value = true
}

const customRequest = ({
    file,
    data,
    headers,
    onFinish,
    onError,
    onProgress
}: UploadCustomRequestOptions) => {
    const formData = new FormData();
    if (data) {
        Object.keys(data).forEach((key) => {
            formData.append(
                key,
                data[key as keyof UploadCustomRequestOptions['data']]
            );
        });
    }
    if (file.file) { // 添加额外的检查
        formData.append("file", file.file as File);
    }
    onProgress({
        percent: 0
    })
    axios.post(import.meta.env.VITE_BASE_API_URL + "/api/file/upload", formData, {
        headers: {
            ...headers,
            'Content-Type': 'multipart/form-data'
        },
        withCredentials: true
    }).then(async (response) => {
        if (response.data.ok) {
            let ava_url = import.meta.env.VITE_BASE_API_URL + response.data.data.filePath + response.data.data.fileName
            fileList.value[0].url = ava_url
            userInfo.value.UserInfo.HeaderField = ava_url
            await UserHeaderFieldUpdate({
                "HeaderField": ava_url,
            })
        }
        onFinish()
    }).catch((error) => {
        onError()
    });
};

</script>
<template>
    <div class="page-cont">
        <div class="user-info-cont">
            <div class="user-info-head-cont">
                <n-avatar round :size="80" :src="userInfo.UserInfo.HeaderField" style="cursor: pointer;"
                    @click="showChangeAvatarModal"></n-avatar>
                <div class="user-nick">
                    {{ userInfo.UserInfo.NickName }}
                </div>
            </div>
            <div>
                上次登录 {{ dateDiff(userInfo.UserInfo.LastLogin) }}
            </div>
            <div style="display: flex; justify-content: space-between; width: 80%;">
                <div class="user-info-item">
                    <div class="user-info-item-number">
                        {{ userInfo.TotalCount }}
                    </div>
                    <div class="user-info-item-name">
                        浏览量
                    </div>
                </div>
                <div class="user-info-item user-cursor" @click="goFollowing(userInfo.UserInfo.ID)">
                    <div class="user-info-item-number">
                        {{ userInfo.FollowingNumber }}
                    </div>
                    <div class="user-info-item-name">
                        关注
                    </div>
                </div>
                <div class="user-info-item user-cursor" @click="goFollower(userInfo.UserInfo.ID)">
                    <div class="user-info-item-number">
                        {{ userInfo.FollowerNumber }}
                    </div>
                    <div class="user-info-item-name">
                        粉丝
                    </div>
                </div>
                <div class="user-info-item user-cursor" @click="goCollections(userInfo.UserInfo.ID)">
                    <div class="user-info-item-number">
                        {{ userInfo.CollectionNumber }}
                    </div>
                    <div class="user-info-item-name">
                        收藏
                    </div>
                </div>
                <div class="user-info-item user-cursor" @click="goArticleList(userInfo.UserInfo.ID)">
                    <div class="user-info-item-number">
                        {{ userInfo.ArticleNumber }}
                    </div>
                    <div class="user-info-item-name">
                        文章
                    </div>
                </div>
            </div>
            <div class="user-follow-btn" @click="showModal = !showModal">
                管理
            </div>
            <n-button class="user-logout" @click="doLogout">
                <n-icon>
                    <LogoutOutlined />
                </n-icon>
            </n-button>
        </div>
        <h3>
            文章管理
        </h3>
        <div>
            <div class="user-article-list-item" v-for="article in userInfo.ArticleList">
                <div class="user-article-list-item-left">
                    <div class="user-article-list-item-title">
                        {{ article.Title }}
                    </div>
                    <div class="user-article-list-item-time">
                        创建于 {{ dateToString(article.CreatedAt) }}
                    </div>
                </div>
                <div class="user-article-list-item-right" @click="goEditArticle(1, article.ArticleID)">
                    <n-icon size="30">
                        <KeyboardArrowRightFilled />
                    </n-icon>
                </div>
            </div>
        </div>
        <div class="user-fix-add-btn" @click="goEditArticle(2, 0)">
            <n-icon size="36">
                <Add />
            </n-icon>
        </div>

    </div>

    <n-modal v-model:show="showModal">
        <n-card style="width: 600px" title="管理" :bordered="false">
            <n-input class="user-modal-input" placeholder="用户名" v-model:value="userInfo.UserInfo.UserName" disabled>
                <template #prefix>
                    <span class="user-modal-label">
                        用户名
                    </span>
                </template>
            </n-input>
            <n-input class="user-modal-input" placeholder="昵称" v-model:value="userInfo.UserInfo.NickName">
                <template #prefix>
                    <span class="user-modal-label">
                        昵称
                    </span>
                </template>
            </n-input>
            <n-input class="user-modal-input" placeholder="邮箱" v-model:value="userInfo.UserInfo.Email">
                <template #prefix>
                    <span class="user-modal-label">
                        邮箱
                    </span>
                </template>
            </n-input>
            <n-input-number class="user-modal-input" placeholder="电话号码" v-model:value="userInfo.UserInfo.PhoneNum"
                :show-button="false">
                <template #prefix>
                    <span class="user-modal-label">
                        电话号码
                    </span>
                </template>
            </n-input-number>
            <n-button @click="update">
                提交
            </n-button>
        </n-card>
    </n-modal>

    <n-modal v-model:show="showAvatarModal">
        <n-card style="width: 130px" title="修改头像" :bordered="false">
            <div style="display: flex; align-items: center;justify-content: center;">
                <n-upload :custom-request="customRequest" :default-file-list="fileList" :multiple="false"
                    list-type="image" :show-file-list="false">
                    <n-avatar round :size="80" :src="fileList[0].url as string" style="cursor: pointer;">
                    </n-avatar>
                </n-upload>
            </div>
        </n-card>
    </n-modal>
</template>
<style>
.user-fix-add-btn {
    z-index: 2;
    position: fixed;
    bottom: 100px;
    left: 50%;
    transform: translateX(-50%);
    border-radius: 50%;
    background: var(--add-btn-background);
    backdrop-filter: blur(2px);
    width: 50px;
    height: 50px;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: 0.3s;
    cursor: pointer;
}

.user-logout {
    position: absolute;
    right: 86px;
    top: 15px;
}

.user-fix-add-btn:hover {
    box-shadow: 0 0 12px var(--second-highlight-color);
    color: var(--color-rev);
    background: var(--second-highlight-color);
}

.user-article-list-item {
    margin: 10px 0;
    height: 30px;
    padding: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: var(--card--background);
    border-radius: 12px;
    transition: 0.3s;
}

.user-article-list-item:hover {
    box-shadow: 0 0 4px var(--second-highlight-color);
}

/* .user-article-list-item-left {} */

.user-article-list-item-title {
    font-weight: bold;
    font-size: large;
}

.user-article-list-item-time {
    font-weight: 200;
    font-size: small;
}

.user-article-list-item-right {
    background: var(--layout--bottom-background);
    border-radius: 50%;
    display: flex;
    justify-content: center;
    cursor: pointer;
    transition: 0.3s;
}

.user-article-list-item-right:hover {
    box-shadow: 0 0 12px var(--second-highlight-color);
    color: var(--color-rev);
    background: var(--second-highlight-color);
}

.user-modal-label {
    font-weight: 200;
    color: var(color);
}

.user-modal-input {
    margin: 6px 0;
}

.n-modal-scroll-content {
    justify-content: center;
}
</style>