import React, {useState} from 'react';
import './ProfileContent.css';
import 'boxicons/css/boxicons.min.css';
import PersonInfo from '../PersonInfo/PersonInfo';
import ChangePassword from '../ChangePassword/ChangePassword';

function ProfileContent(props) {
    const [currentIndex, setCurrentIndex] = useState(0);

    return (
        <div className="profile-container">
            <div className="profile-sidebar">
                <div className="profile-user">
                    <div className="profile-userpic">
                        <img src="https://scontent.fhan14-4.fna.fbcdn.net/v/t1.6435-9/182741929_107870648128506_6644574165595211881_n.jpg?_nc_cat=102&ccb=1-7&_nc_sid=09cbfe&_nc_ohc=5kAhAr4PEkoAX-zfreP&_nc_ht=scontent.fhan14-4.fna&oh=00_AfCJqeKHCh0PvtaZ9Y9AiY_0vUc6X8QkOtlc9loxucFV0w&oe=649DFFEE" width="100" height="100" alt="" />
                        <div className="round">
                        <input type="file" />
                        <i class='bx bx-camera'style={{ color: '#fff' }} ></i>
                        </div>
                    </div>
                    <div className="profile-usertitle">
                        <div className="profile-usertitle-name">
                            Hoàn Bằng
                        </div>
                        <div className="profile-usertitle-role">
                            Nhân viên
                        </div>
                    </div>
                </div>

                <div className='profile-userinfo'>
                <div className='profile-userinfo-item'>
                        <div>
                            <i class='bx bx-calendar-alt' ></i> Ngày sinh:
                        </div>
                        <div>
                            15/11/2002
                        </div>
                    </div>
                    
                    <div className='profile-userinfo-item'>
                        <div>
                            <i class='bx bx-user'></i> Giới tính:
                        </div>
                        <div>
                            Nam
                        </div>
                    </div>
                    
                    <div className='profile-userinfo-item'>
                        <div>
                            <i class='bx bx-phone' ></i> Số điện thoại:
                        </div>
                        <div>
                            0388586955
                        </div>
                    </div>
                    
                    <div className='profile-userinfo-item'>
                        <div>
                            <i class='bx bx-envelope' ></i> Email:
                        </div>
                        <div>
                            hoannk1511@gmail.com
                        </div>
                    </div>
                </div>
                
                <div className="profile-usermenu">
                    <ul className="profile-usermenu-nav">
                        <li className={currentIndex === 0 ? "active" : ""}>
                            <a href="#" onClick={() => setCurrentIndex(0)}>
                            Cập nhật thông tin </a>
                        </li>
                        <li className={currentIndex === 1 ? "active" : ""}>
                            <a href="#" onClick={() => setCurrentIndex(1)}>
                            Thay đổi mật khẩu </a>
                        </li>
                    </ul>
                </div>
                        
                </div>
                <div className="profile-content">
                    {currentIndex === 0 ? <PersonInfo/> : <ChangePassword/>}
                </div>
        </div>
    );
}

export default ProfileContent;